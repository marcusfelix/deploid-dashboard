# Setup

In oder to deploy this app you will need create a Firebase Project

`docker build -t caster .`

`docker run -dp 8090:8090 caster`

`flutter pub run icons_launcher:create`

`flutter pub global activate rename`

`flutter pub global run rename --bundleId com.deploid.app`

`flutter pub global run rename --appname "App name"`