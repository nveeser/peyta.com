#!/usr/bin/python

import subprocess
import sys
import os
import argparse
import time
import datetime
import shutil
import re
import json
import fcntl
import traceback

from os.path import join
from subprocess import Popen, PIPE, STDOUT
from contextlib import contextmanager

UNRAR='/usr/bin/urar'
UNZIP='/usr/bin/unzip'
FFMPEG='/home/nicholas/bin/ffmpeg'
FFPROBE='/usr/bin/ffprobe'

FFMPEG_PROFILES = { 
    "core" :  [ 
        '-y',
        '-codec:v', 'libx264',
        '-profile:v', 'high', 
        '-preset', 'slow',
        '-b:v', '9.8M',
        '-crf', '21',
        '-pix_fmt', 'yuv420p',
        '-codec:a', 'libfdk_aac',
        '-b:a', '128k' 
    ],
    "ipad" : [
        '-y',
        '-codec:a', 'libfdk_aac',
        '-b:a', '160k', 
        '-ac', '2', 
        '-strict',  'experimental',
        '-s', 'hd1080',
        '-codec:v', 'libx264',
        '-b:v', '1200k', 
        '-profile:v', 'high',
        '-level', '4.2',
        '-preset', 'slow',
        '-pix_fmt', 'yuv420p',
        '-maxrate', '10000000',
        '-bufsize', '10000000', 
        '-f', 'mp4',
    ],
    "roku" : [
        '-y',
        '-aspect', '16:9',
        '-vf', 'yadif=0:1',
        '-codec:v', 'libx264',
        '-bufsize', '2000k',
        '-profile:v', 'high',
        '-level', '4.0',
        '-s', '720x480',
        '-b:v', '1000k',
        '-codec:a', 'libfdk_aac',
        '-b:a', '128k',
        '-vpre', 'roku',
        '-r', '29.97',
        '-adrift_threshold', '0.01',
    ],        
        
}

FFPROBE_ARGS = [  
    '-print_format', 'json',
    '-loglevel', '0',
    '-show_streams'
]

EXT_VIDEOS = [ 'wmv', 'mov', 'mp4', 'avi', 'flv', 'm4v', 'mpeg', 'asf', 'mkv' ]
EXT_ARCHIVES = ['zip', 'rar']
EXT_PICTURES = [ 'jpg', 'png', 'jpeg', 'txt' ]
EXT_TEXT = [ 'txt', 'nfo' ]

TYPE_MAP = {
    'pics' : EXT_ARCHIVES + EXT_PICTURES + EXT_TEXT,
    'vids' : EXT_VIDEOS,
    'all' : EXT_ARCHIVES + EXT_PICTURES + EXT_TEXT + EXT_VIDEOS,
}

MAX_BITRATE=10.0

PROGRESS_DIR=".progress"

def build_parser():
    parser = argparse.ArgumentParser(prog='ghost',
                                     description='Process some files')

    parser.add_argument('-n', '--testing', dest='testing', action='store_true',
                        help='Set the testing flag: dont run commands')
    parser.add_argument('--verbose', '-v', action='count')
    parser.add_argument('--lockdir', required=False, default='/tmp/ghost_locks')
    parser.add_argument('--profile', required=False, default='core')

    sub = parser.add_subparsers(help='sub-command help', dest='command')
    p1 = sub.add_parser("allextract",
                        help='extract all archives from the current directory')


    p2 = sub.add_parser("transcode",
                        help='transcode video file, outputting to the specififed directory')
    p2.add_argument("infile")
    p2.add_argument("outdir")


    p3 = sub.add_parser("process",
                        help='Process all files in the specfified directory')
    p3.add_argument("source")
    p3.add_argument("dest")
    p3.add_argument("--source-base", required=False, default=None)
    p3.add_argument("--extensions", required=False, default='all', choices=['vids', 'pics', 'all'])

    p4 = sub.add_parser("bulk",
                        help='Bulk process files, moving original to the backup')
    p4.add_argument("source")
    p4.add_argument("--bulk-root", required=False, default="/m2/bag")
    p4.add_argument("--backup-root", required=False, default="/scratch/backup/movies")

    p4 = sub.add_parser("rsort",
                       help='Sort files into directories')
    p4.add_argument("source")

    return parser


class Processor:
    def __init__(self, source, dest, source_base):
        self.source = os.path.normpath(source)
        self.dest = os.path.normpath(dest)
        if source_base:
            source_base = os.path.normpath(source_base)
        self.source_base = Processor.make_source_base(source, source_base)
        self.progress_dir = os.path.join(dest, PROGRESS_DIR)

    @staticmethod
    def make_source_base(source, source_base):
        if not source_base:
            return source

        if not (source == source_base or source.startswith(source_base)):
            raise Exception("source_base (%s) must be a prefix of the input dir (%s) " 
                            % (source_base, source))
        return source_base


class TorrentFile:
    def __init__(self, filename, processor):
        self.filename = filename
        self.p = processor

    def has_pfile(self):
        return os.path.exists(self.progress_file())

    def progress_file(self):
        return self.filename.replace(self.p.source_base, self.p.progress_dir)
        
    def extension(self):
        return _extension(self.filename)

    def outdir(self):
        return os.path.dirname(self.filename).replace(self.p.source, self.p.dest)

    def basename_outdir(self):
        """ /path/source/file.zip -> /path/output/file"""
        basename, _ = os.path.splitext(self.filename)
        return basename.replace(self.p.source, self.p.dest)

    def derived_outdir(self):
        if self.extension() in EXT_ARCHIVES:
            return self.basename_outdir()
        else:
            return self.outdir()
    
    def __repr__(self):
        return ("Filename: %s\n Source:%s\n Dest:%s\n SourceBase:%s" % 
                (self.filename, self.p.source, self.p.dest, self.p.source_base))


        
class Ghost:
    def __init__(self, profile, verbose=0, testing=False):
        self.verbose = verbose
        self.testing = testing
        self.profile = profile
        self._prefix = ""

    def allextract(self):
        cwd = os.getcwd()
        archives = _yield_files_by_ext(cwd, ARCHIVES)

        for filepath in archives:
            basename, _ = os.path.splitext(filepath)
            if not os.path.isdir(basename):
                self.extract(filepath, basename)


    def bulk_transcode(self, source, bulk_root, backup_root):
        source = os.path.abspath(source)

        def not_mp4(file):
            if _extension(file) != 'mp4':
                return True

            self.log(1, "IGNORE MP4: " + file)
            return False

        files = _yield_files_by_ext(source, EXT_VIDEOS)
        files = [ f for f in files if not_mp4(f) ]
        files = self.prefix(files)

        for inpath in files:
            if not inpath.startswith(bulk_root):
                raise Exception("Path is not within --bulk-root=%s: %s" % (bulk_root, inpath))

            backup = inpath.replace(bulk_root, backup_root)
            outdir = os.path.dirname(inpath)

            try:
                self.transcode(inpath, outdir)
                self._backup(inpath, backup)

            except Exception as e:
                self.log(0, "Error with file: %s" % e)



    def process_dir(self, source, dest, source_base, extensions):
        p = Processor(source, dest, source_base)

        files = _yield_files_by_ext(source, extensions)
        files = [ TorrentFile(f, p) for f in files ]
        files = sorted(files, key=lambda t: t.filename)

        def tfile_filter(t_file):
            if (t_file.extension() not in EXT_ARCHIVES 
                and os.path.exists(t_file.outdir())
                and not os.path.isdir(t_file.outdir())):
                #print("BROKEN: %s -> %s" % (t_file.filename, t_file.outdir()))
                #os.unlink(t_file.progress_file())
                return True

            if t_file.has_pfile():
                self.log(1, "SKIP: ", t_file.filename)
                return False
            return True

        #self.log(1, "FILE: ", e)
        files = [ f for f in files if tfile_filter(f) ]
        files = self.prefix(files)

        for t_file in files:
            try:
                self.log(1, "FILE: ", t_file.filename)
                if t_file.extension() in EXT_VIDEOS:
                    self.transcode(t_file.filename, t_file.outdir())

                elif t_file.extension() in EXT_ARCHIVES:
                    self.extract(t_file.filename, t_file.basename_outdir())

                elif t_file.extension() in EXT_PICTURES + EXT_TEXT:
                    self._copyfile(t_file.filename, t_file.outdir())

                # Touch the Pfile
                self._progress_file(t_file.progress_file())

            except Exception as e:
                self.log(0, "Error with file: %s" % e)
                traceback.print_exc()


    def transcode(self, source, destdir):
        """
        source: "/m2/bag/foo.wmv"
        destdir: "/m2/bag/output"

        outfile: "/m2/bag/output/foo.mp4
        Args:
        source: path to the file to transcode
        destdir: path to the directory to write to
        """

        if not os.path.exists(destdir):
            os.makedirs(destdir)

        basename = os.path.basename(source)

        # if _extension(source) == "mp4" and not is_high_bitrate_mp4(source):            
        #     self._copyfile(source, destdir)
        #     return 

        base, _ = os.path.splitext(basename)
        dest = join(destdir, base + ".mp4")
        args = [ FFMPEG, '-i', source ] + FFMPEG_PROFILES[self.profile] + [ dest ]

        self.log(1, "TRANSCODE: %s -> %s" % (source, dest))

        code = self._run(args)
        
        if code != 0:
            raise Exception("Non-zero return code (%d) for %s" % (code, " ".join(args)))


    def extract(self, source, destdir):
        if os.path.exists(destdir):
            self.log(1, "EXISTS: ", destdir)
            return

        self.log(1, "EXTRACT: %s -> %s" % (source, destdir))

        tempdest = "%s.TEMP" % destdir
        if not os.path.exists(tempdest) and not self.testing:
            os.makedirs(tempdest)

        try:
            ext = _extension(source)
            if ext == 'zip':
                args = [UNZIP, '-o', '-j', '-d', tempdest, source]

            elif ext == 'rar':
                args = ['7z', 'x', '-o', tempdest, source]

            code = self._run(args)

            if os.path.isdir(tempdest) and not self.testing:
                os.rename(tempdest, destdir)

            if code not in [ 0, 2 ]:
                raise Exception("Non-zero return code (%d) for %s" % (code, " ".join(args)))


        finally:
            if os.path.isdir(tempdest):
                self.log(1, "CLEANUP: ", tempdest)
                shutil.rmtree(tempdest, True)


    def regex_sort(self, source):
        source = os.path.abspath(source)
        for inpath in _yield_files_by_ext(source, ['jpg']):
            basepath = os.path.dirname(inpath)
            filename = os.path.basename(inpath)
            m = re.match("((.*\w)_|(.*[a-zA-Z]))\d+\.\w{3}$", filename)
            if not m:
                self.log(1, "SKIP: ", filename)
                continue

            dirname = m.group(2) or m.group(3)
            outpath = os.path.join(basepath, dirname, filename)
            self.log(0, "MOVE: %s -> %s" % (inpath, outpath))

            if self.testing:
                continue

            if not os.path.isdir(os.path.dirname(outpath)):
                os.makedirs(os.path.dirname(outpath))

            shutil.move(inpath, outpath)


    def log(self, level, *args):
        if self.verbose is not None and self.verbose >= level:
            print(self._prefix)
            print("".join(args))


    def prefix(self, files):
        files = list(files)
        total = len(files)
        for num, f in enumerate(files):
            self._prefix = "[%d/%d] " % (num+1, total)
            yield f


    def _make_base(self, path):
        if self.testing:
            return

        dirname = os.path.dirname(path)
        if not os.path.exists(dirname):
            os.makedirs(dirname)


    def _run(self, args):
        start = time.time()
        self.log(2, "RUN: ", " ".join(args))

        if self.testing:
            return 0

        proc = Popen(args, stdout=PIPE, stderr=STDOUT)
        while True:
            line = proc.stdout.readline().decode('utf-8')
            if line == '':
                break

            self.log(2, ">> ", line.rstrip())

        proc.stdout.close()
        proc.wait()

        if self.verbose:
            elapsed = datetime.timedelta(0, time.time() - start)
            self.log(2, "DURATION: ", str(elapsed))

        return proc.returncode


    def _copyfile(self, source, destdir):
        self.log(1, "COPY: %s -> %s" % (source, destdir))

        if self.testing:
            return

        if os.path.isfile(destdir):
            self.log(0, "CLEANUP: %s" % destdir)
            os.unlink(destdir)

        if not os.path.isdir(destdir):    
            os.makedirs(destdir)

        shutil.copy2(source, destdir)
        

    def _backup(self, inpath, outpath):
        self.log(1, "BACKUP: %s -> %s" % (inpath, outpath))

        if self.testing:
            return

        outdir = os.path.dirname(outpath)
        if not os.path.isdir(outdir):
            os.makedirs(outdir)

        shutil.move(inpath, outpath)


    def _progress_file(self, pfile):
        pdir = os.path.dirname(pfile)
        if not os.path.isdir(pdir):
            os.makedirs(pdir)
        if not self.testing:
            with open(pfile, 'w'):
                os.utime(pfile, None)



def is_high_bitrate_mp4(source):
    data = _probe_file(source) 
    bitrate = None

    for stream in data.get('streams', {}):
        tag_str = "TAG(%s)" % stream['codec_tag']
        codec = stream.get('codec_name', tag_str)
        if codec == 'h264':
            bitrate = int(stream.get('bit_rate', '0'))
            bitrate = float(bitrate) / 1024 / 1024            
            break

        
    if bitrate and bitrate > MAX_BITRATE:
        print("H264: %0.2fMib/s %s" % (bitrate, source))
        return True
    return False


def _probe_file(source):
    args = [ FFPROBE, source ] + FFPROBE_ARGS 
    p = subprocess.Popen(args, stdin=subprocess.PIPE, stdout=subprocess.PIPE, stderr=subprocess.PIPE)    
    (out, err) = p.communicate()
    return json.loads(unicode(out, errors='replace'))


def _duration(s):
    hours, remainder = divmod(s, 3600)
    minutes, seconds = divmod(remainder, 60)
    return '%dh:%dm:%ds' % (hours, minutes, seconds)


def _extension(filepath):
    """Returns extension, removing the leading '.' and lowercase"""
    base, ext = os.path.splitext(filepath)
    return ext[1:].lower()


def _yield_files_by_ext(path, extensions):
    for root, _, files in os.walk(path):
        for name in files:
            filepath = join(root, name)
            _, ext = os.path.splitext(filepath)
            ext = ext[1:].lower() 
            if ext in extensions:
                yield filepath


@contextmanager
def file_locks(lock_files):
    for lock in lock_files:
        if os.path.exists(lock):
            sys.exit(0)
 
    try:
        for lock in lock_files:
            f = open(lock, 'w')
            fcntl.lockf(f.fileno(), fcntl.LOCK_EX | fcntl.LOCK_NB)
            
        yield

    finally:
        for lock in lock_files:
            if os.path.exists(lock):
                os.remove(lock)



if __name__ == "__main__":
    parser = build_parser()
    args = parser.parse_args(sys.argv[1:])

    ghost = Ghost(args.profile, verbose=args.verbose, testing=args.testing)

    if args.command == "allextract":
        ghost.allextract()

    elif args.command == "transcode":
        ghost.transcode(args.infile, args.outdir)

    elif args.command == "bulk":
        ghost.bulk_transcode(args.source, args.bulk_root, args.backup_root)

    elif args.command == "rsort":
        ghost.regex_sort(args.source)

    elif args.command == "process":
        if not os.path.isdir(args.lockdir):
            os.makedirs(args.lockdir)
        ext_list = TYPE_MAP[args.extensions]
        locks = [ os.path.join(args.lockdir, ext) for ext in ext_list ]
        with file_locks(locks):
            ghost.process_dir(args.source, args.dest, args.source_base, ext_list)

    else:
        print("Unknown Command: %s" % command)
        sys.exit(1)
