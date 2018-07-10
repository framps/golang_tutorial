# Sample multiuser golang webserver which provides the famous highlow game and can be run locally and deployed on Bluemix

This is a small sample webserver implemented in golang which provides the famous highlow game and can be used by multiple users in parallel.

There exist various make targets (See *Makefile* for details). The most interesting targets are

1. ```make run``` to start the webserver on the local system. Open your browser and connect to http://127.0.0.1:8080 to play the game on your local system and watch all logging messages in realtime.
2. ```make deploy``` to deploy the webserver on Bluemix. You have to have an existing Bluemix account. If you don't have one then register for a evaluation account which can be used for 30 day for free. 

 At the end the url to access the deployed sample application will be displayed. It will start with *highLowGameServerMultiuser* and have some random generated trailer as *-divisional-squattocracy* for example thus the final access url will be for this example http://highLowGameServerMultiuser-divisional-squattocracy

Enjoy playing highlow [here](https://highlowgameservermultiuser.mybluemix.net/)
