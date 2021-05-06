with open("conos_urls.txt", "r") as f:
    lines = f.readlines()
with open("conos_urls.txt", "w") as f:
    for line in lines:
        if 'flickr' in line.strip("\n"):
            f.write(line)
