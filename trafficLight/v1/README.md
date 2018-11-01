# traffic light sample program configuration for two trafficlights

Two trafficlights are simulated.
* trafficlight1: GPIO2-GPIO4 LEDs red, yellow and green
* trafficlight2: GPIO5-GPIO7 LEDs red, yellow and green

Trafficlight programs are executed in sequence

## Trafficlight configuration

### Define GPIO to Pin mapping

Config file GPIO.json:

```
[
17, 18,
27, 22, 23,
24, 25, 4
]
```

```
17 - GPIO0 (unused)
18 - GPIO1 (unused)
27 - GPIO2 is BCM pin 27 which is trafficlight1 red LED
22 - GPIO3 is BCM pin 22 which is trafficlight1 yellow LED
...
25 - GPIO6 is BCM pin 25 which is trafficlight2 yellow LED
4 - GPIO7 is BCM pin 4 which is trafficlight2 green LED
```

### Traffic light program definitions

Config file: programs.json

Default definitions:
```
{
   "Normal1": {
      "name": "Normal1",
      "phases": [
         {
            "light": 1,
            "ticks": 3
         },
         {
            "light": 2,
            "ticks": 2
         },
         {
            "light": 3,
            "ticks": 4
         },
         {
            "light": 4,
            "ticks": 1
         }
      ]
   },
   "Normal2": {
      "name": "Normal2",
      "phases": [
         {
            "light": 1,
            "ticks": 4
         },
         {
            "light": 2,
            "ticks": 1
         },
         {
            "light": 3,
            "ticks": 4
         },
         {
            "light": 4,
            "ticks": 1
         }
      ]
   },
   "Normal3": {
      "name": "Normal3",
      "phases": [
         {
            "light": 1,
            "ticks": 4
         },
         {
            "light": 2,
            "ticks": 1
         },
         {
            "light": 2,
            "ticks": 1
         },
         {
            "light": 3,
            "ticks": 4
         },
         {
            "light": 3,
            "ticks": 1
         },
         {
            "light": 4,
            "ticks": 1
         }
      ]
   }
}
```

## Options
```
-debug
     Write debug messages
-leds
     Drive LEDs
-monitor
     Monitor LEDs on screen (default true)
```
