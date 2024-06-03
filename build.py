	
import sys
import os
import subprocess

def main():
	files = list(filter(lambda name: name.endswith('.go'), os.listdir('.')))
	files = list(filter(lambda name: name not in map(lambda arg: arg.split('=')[1], filter(lambda arg: arg.startswith('--build-ignore='), sys.argv[1:])), files))
	rc = subprocess.call(['go', 'build', '-o', 'main.exe']+files)
	if rc == 0:
		subprocess.call(['main.exe']+sys.argv[1:])

if __name__ == '__main__':
	main()