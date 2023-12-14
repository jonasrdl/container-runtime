# container-runtime

To setup groups, permissions and create the main folder, run:
```
sudo ./install.sh
```

Pull images using
```
./pull <image-name>
```

Then you can run the container with a desired image using
```
./container-runtime run <image-name> [command]
```
Command is optional, if you don't provide it, the default command from the image will be used.

You can list all containers using
```
./container-runtime list
```

You can delete a container using
```
./container-runtime delete <container-id>
```

You can also delete all containers using
```
./container-runtime delete --all
```