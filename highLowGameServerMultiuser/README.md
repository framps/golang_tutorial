# Sample high low multiuser golang webserber

This is a small sample webserver implemented in golang which provides the famous high-low game and can be used by multiple users in parallel.

There exist various make targets. The most interesting targets are

1. ```make run``` to start the webserver on the local system. Open your browser and connect to http://127.0.0.1:8080 to play the game and watch all logging messages in realtime.
2. ```make deploy``` to deply the webserver on Bluemix. At the end ```cf apps``` will display the url to access the deployed sample application. It wil start with *highLowGameServerMultiuser* and have some random trailer like *-divisional-squattocracy* for example thus the final access url will be http://highLowGameServerMultiuser-divisional-squattocracy

Enjoy playing highlow :-)
