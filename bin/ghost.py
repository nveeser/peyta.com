#!/usr/bin/python

import subprocess
import sys
import os
import argparse
import time
import datetime
import shutil
import re

from os.path import join
from subprocess import Popen, PIPE, STDOUT
from contextlib import contextmanager

UNRAR='/usr/bin/unrar'
UNZIP='/usr/bin/unzip'
FFMPEG='/home/nicholas/bin/ffmpeg'

FFMPEG_ARGS= [ '-y',
               '-c:v', 'libx264',
               '-preset', 'slow',
               '-crf', '21',
               '-pix_fmt', 'yuv420p',
               '-c:a', 'libfdk_aac',
               '-b:a', '128k' ]

EXT_ARCHIVES = ['zip', 'rar']
EXT_VIDEOS = [ 'wmv', 'mov', 'mp4', 'avi', 'flv', 'm4v' ]
EXT_ALL = EXT_ARCHIVES + EXT_VIDEOS

PROGRESS_DIR=".progress"

def build_parser():
    parser = argparse.ArgumentParser(prog='ghost',
                                     description='Process some files')

    parser.add_argument('-n', '--testing', dest='testing', action='store_true',
                        help='Set the testing flag: dont run commands')
    parser.add_argument('--verbose', '-v', action='count')

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
    p3.add_argument("--extensions", required=False, default=None, choices=['vids', 'pics'])

    p4 = sub.add_parser("bulk",
                       help='Bulk process files, moving original to the backup')
    p4.add_argument("source")
    p4.add_argument("--bulk-root", required=False, default="/m2/bag")
    p4.add_argument("--backup-root", required=False, default="/m2/bag/.Backup")

    p4 = sub.add_parser("rsort",
                       help='Sort files into directories')
    p4.add_argument("source")

    return parser


class Ghost:
    def __init__(self, verbose=0, testing=False):
        self.verbose = verbose
        self.testing = testing
        self._prefix = ""

    def allextract(self):
        cwd = os.getcwd()
        archives = _yield_files(cwd, ARCHIVES)

        for filepath in archives:
            basename, _ = os.path.splitext(filepath)
            if not os.path.isdir(basename):
                self.extract(filepath, basename)


    def regex_sort(self, source):
        source = os.path.abspath(source)
        for inpath in _yield_files(source, ['jpg']):
            basepath = os.path.dirname(inpath)
            filename = os.path.basename(inpath)
            m = re.match("((.*\w)_|(.*[a-zA-Z]))\d+\.\w{3}$", filename)
            if not m:
                self.log("SKIP: ", filename)
                continue

            dirname = m.group(2) or m.group(3)
            outpath = os.path.join(basepath, dirname, filename)
            self.log("MOVE: %s -> %s" % (inpath, outpath))

            if self.testing:
                continue

            if not os.path.isdir(os.path.dirname(outpath)):
                os.makedirs(os.path.dirname(outpath))

            shutil.move(inpath, outpath)


    def bulk_transcode(self, source, bulk_root, backup_root):
        source = os.path.abspath(source)

        def not_mp4(file):
            if _extension(file) != 'mp4':
                return True

            if self.verbose:
                self.log("IGNORE MP4: " + file)

            return False

        files = _yield_files(source, EXT_VIDEOS)
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
                self.log("Error with file: %s" % e)


    def process(self, indir, dest, base, extensions=EXT_ALL):
        if base:
            if not indir.startswith(base) and not indir == base:
                raise Exception("Source dir (%s) must be a prefix of the input dir (%s) " % (base, indir))
            source = base
        else:
            source = indir

        progress_dir = os.path.join(dest, PROGRESS_DIR)

        def has_pfile(file):
            pfile = file.replace(source, progress_dir)
            if not os.path.exists(pfile):
                return True

            if self.verbose:
                self.log("SKIP: ", file)

            return False

        files = _yield_files(indir, extensions)
        files = [ f for f in files if has_pfile(f) ]
        files = self.prefix(files)

        for inpath in files:
            pfile = inpath.replace(source, progress_dir)

            if _extension(inpath) in EXT_VIDEOS:
                outdir = os.path.dirname(inpath).replace(source, dest)
                self.transcode(inpath, outdir)

            elif _extension(inpath) in EXT_ARCHIVES:
                basename, _ = os.path.splitext(inpath)
                outdir = basename.replace(source, dest)
                self.extract(inpath, outdir)

            # Touch the Pfile
            self._progress_file(pfile)


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
        base, _ = os.path.splitext(basename)
        dest = join(destdir, base + ".mp4")
        args = [ FFMPEG, '-i', source ] + FFMPEG_ARGS + [ dest ]

        if self.verbose:
            self.log("TRANSCODE: %s -> %s" % (source, dest))

        code = self._run(args)

        if code != 0:
            raise Exception("Non-zero return code (%d) for %s" % (code, " ".join(args)))


    def extract(self, source, destdir):
        if os.path.exists(destdir):
            self.log("EXISTS: ", destdir)
            return

        if self.verbose:
            self.log("EXTRACT: %s -> %s" % (source, destdir))

        basename, _ = os.path.splitext(source)
        if os.path.isdir(basename):
            if self.verbose:
                self.log("MOVE EXISTING [%s] -> [%s]" % (basename, destdir))

            self._make_base(destdir)
            os.rename(basename, destdir)
            return


        tempdest = "%s.TEMP" % destdir
        if not os.path.exists(tempdest) and not self.testing:
            os.makedirs(tempdest)

        try:
            ext = _extension(source)
            if ext == 'zip':
                args = [UNZIP, '-o', '-j', '-d', tempdest, source]

            elif ext == 'rar':
                args = [UNRAR, 'e', '-y', source, tempdest]

            code = self._run(args)

            if os.path.isdir(tempdest) and not self.testing:
                os.rename(tempdest, destdir)

            if code not in [ 0, 2 ]:
                raise Exception("Non-zero return code (%d) for %s" % (code, " ".join(args)))


        finally:
            if os.path.isdir(tempdest):
                if self.verbose:
                    self.log("CLEANUP: ", tempdest)
                shutil.rmtree(tempdest, True)


    def log(self, *args):
        print self._prefix,
        print "".join(args)


    def prefix(self, files):
        files = list(files)
        total = len(files)
        for num, inpath in enumerate(files):
            self._prefix = "[%d/%d] " % (num+1, total)
            yield inpath


    def _make_base(self, path):
        if self.testing:
            return

        dirname = os.path.dirname(path)
        if not os.path.exists(dirname):
            os.makedirs(dirname)


    def _run(self, args):
        start = time.time()
        if self.verbose > 1:
            self.log("RUN: ", " ".join(args))

        if self.testing:
            return 0

        proc = Popen(args, stdout=PIPE, stderr=STDOUT)
        while True:
            line = proc.stdout.readline()
            if line == '':
                break

            if self.verbose > 1:
                self.log(">> ", line.rstrip())

        proc.stdout.close()
        proc.wait()

        if self.verbose:
            elapsed = datetime.timedelta(0, time.time() - start)
            self.log("Duration: ", str(elapsed))

        return proc.returncode


    def _backup(self, inpath, outpath):
        if self.verbose:
            self.log("BACKUP: %s -> %s" % (inpath, outpath))

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



def _duration(s):
    hours, remainder = divmod(s, 3600)
    minutes, seconds = divmod(remainder, 60)
    return '%dh:%dm:%ds' % (hours, minutes, seconds)


def _extension(filepath):
    """Returns extension, removing the leading '.' and lowercase"""
    base, ext = os.path.splitext(filepath)
    return ext[1:].lower()


def _yield_files(path, extensions):
    for root, _, files in os.walk(path):
        for name in files:
            filepath = join(root, name)
            _, ext = os.path.splitext(filepath)
            if ext[1:].lower() in extensions:
                yield filepath


@contextmanager
def file_lock(lock_file):
    if os.path.exists(lock_file):
        sys.exit(0)
    else:
        try:
            f = open(lock_file, 'w').write("1")
            fcntl.lockf(f, fcntl.LOCK_EX | fcntl.LOCK_NB)
            yield
        finally:
            os.remove(lock_file)



if __name__ == "__main__":
    parser = build_parser()
    args = parser.parse_args(sys.argv[1:])

    ghost = Ghost(verbose=args.verbose, testing=args.testing)

    if args.command == "allextract":
        ghost.allextract()

    elif args.command == "transcode":
        ghost.transcode(args.infile, args.outdir)

    elif args.command == "bulk":
        ghost.bulk_transcode(args.source, args.bulk_root, args.backup_root)

    elif args.command == "rsort":
        ghost.regex_sort(args.source)

    elif args.command == "process":
        if args.extensions == 'vids':
            ghost.process(args.source, args.dest, args.source_base, EXT_VIDEOS)

        elif args.extensions == 'pics':
            ghost.process(args.source, args.dest, args.source_base, EXT_ARCHIVES)

        else:
            ghost.process(args.source, args.dest, args.source_base, EXT_ALL)


    else:
        print "Unknown Command: %s" % command
        sys.exit(1)
