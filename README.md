# container-runtime

Pull images using
```
./pull <image-name>
```

Then you can run the container with a desired image using
```
sudo ./container-runtime run <image-name> [command]
```

Command is optional, if you don't provide it, the default command from the image will be used.

You can execute a command inside the container using
```
sudo ./container-runtime exec <container-id>
```

You can list all containers using
```
sudo ./container-runtime list
```

You can delete a container using
```
sudo ./container-runtime delete <container-id>
```

You can also delete all containers using
```
sudo ./container-runtime delete --all
```