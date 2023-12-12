# container-runtime

Pull images using
```
./pull <image-name>
```

Then you can run the container with a desired image using
```
sudo ./container-runtime <image-name> [command]
```

Command is optional, if you don't provide it, the default command from the image will be used.
