import os
import patoolib

path = "C:\\Users\\iphon\\Desktop\\DEMOProject\\All_Rars"
output_dir = "C:\\Users\\iphon\\Desktop\\DEMOProject\\More_Demos"

def getNames(path):
    for root, dirs, filenames in os.walk(path):
        for name in filenames:
            _, extension = os.path.splitext(name)
            if extension.lower() == '.rar':
                file_path = os.path.join(root, name)
                try:
                    patoolib.extract_archive(file_path, outdir=output_dir)
                    os.remove(file_path)
                except Exception as e:
                    print(f"Error extracting {file_path}: {e}")

getNames(path)