# traffic light sample program for two trafficlights

## Abstract

Two trafficlights TL1 and TL2 are simulated (see below their location in the crossing). A traffic light manager loads different traffic light programs in an endless loop, runs them for 15 seconds and switches to the next program. All programs are separated by a 5 second warning blink program which blinks the yellow LEDs. CTRL-C will terminate the endless loop.

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

#### light

| light | green | yellow  | red |
| ----- |:-----:| -------:| --: |
| 0     | off   | off     | off |
| 1     | on    | off     | off |
| 2     | off   | on      | off |
| 3     | off   | off     | on  |
| 4     | off   | on      | on  |

#### ticks

Defines the duration to light LED and is the number of ticks received from traffic manager

#### programs.json

Config file: programs.json

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

## Program options
```
-debug
     Write debug messages
-leds
     Drive LEDs
-monitor
     Monitor LEDs on screen (default true)
```
