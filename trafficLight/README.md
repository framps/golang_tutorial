# traffic light sample program for two trafficlights

## Abstract

Two trafficlights TL1 and TL2 are simulated (see below their location in the crossing). A traffic light manager loads different traffic light programs in an endless loop, runs them for 15 seconds and switches to the next program. All programs are separated by a 5 second warning blink program which blinks the yellow LEDs. CTRL-C will terminate the endless loop.

Either build the executable for Raspberry with `buildAndRun` or just execute `trafficlight_arm`.

### Location of traffic lights

```
         TL1
          |  
      ---   --- TL2
          |
```

## Configuration

Mapping of Raspberry pins to GPIOs and traffic light programs can be configured. A default configuration exists and
1. uses the [BCM GPIO mapping](http://wiringpi.com/pins/)
2. three sample traffic light programs

The mapping from GPIOs to pins is defined in GPIO.json. Traffic light programs are defined in programs.json.

### Mapping of GPIO to LEDs

* TL1: uses GPIO2-GPIO4 LEDs red(2), yellow(3) and green(4)
* TL2: uses GPIO5-GPIO7 LEDs red(5), yellow(6) and green(7)

## Default trafficlight configuration

### GPIO to Pin mappings

#### GPIO.json

```
[
17, 18,
27, 22, 23,
24, 25, 4
]
```

```
17: GPIO0 (unused)
18: GPIO1 (unused)
27: GPIO2 is BCM pin 27 which is trafficlight1 red LED
22: GPIO3 is BCM pin 22 which is trafficlight1 yellow LED
...
25: GPIO6 is BCM pin 25 which is trafficlight2 yellow LED
04: GPIO7 is BCM pin 4 which is trafficlight2 green LED
```

### Traffic light program definitions

Programs are a sequence of tuples which define the LEDs to light and how long they are displayed (number of ticks from traffic manager).

#### lights

| lights | green | yellow  | red |
| ------ |:-----:| -------:| --: |
| 0      |       |         |     |
| 1      | on    |         |     |
| 2      |       | on      |     |
| 3      |       |         | on  |
| 4      |       | on      | on  |

#### ticks

Defines the duration to light LED and is the number of ticks received from traffic manager

#### programs.json

```
{
   "Normal1": {
      "name": "Normal1",
      "phases": [
         {
            "lights": 1,
            "ticks": 3
         },
         {
            "lights": 2,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 3
         },
         {
            "lights": 3,
            "ticks": 3
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 4,
            "ticks": 1
         },
         {
            "lights": 1,
            "ticks": 3
         }
      ],
      "clock_speed": 1000000000
   },
   "Normal2": {
      "name": "Normal2",
      "phases": [
         {
            "lights": 1,
            "ticks": 3
         },
         {
            "lights": 2,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 3
         },
         {
            "lights": 3,
            "ticks": 3
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 4,
            "ticks": 1
         },
         {
            "lights": 1,
            "ticks": 3
         }
      ],
      "clock_speed": 1000000000
   },
   "Normal3": {
      "name": "Normal3",
      "phases": [
         {
            "lights": 1,
            "ticks": 3
         },
         {
            "lights": 2,
            "ticks": 1
         },
         {
            "lights": 2,
            "ticks": 1
         },
         {
            "lights": 3,
            "ticks": 3
         },
         {
            "lights": 3,
            "ticks": 1
         },
         {
            "lights": 4,
            "ticks": 1
         }
      ],
      "clock_speed": 1000000000
   }
}
```

## Program options
```
-leds
     Use LEDs connected to GPIOs (run on Raspberry)
-monitor
     Simulate LEDs on screen (default true) (run on any non Raspberry system)
```
