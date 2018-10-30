# traffic light sample program configuration for two trafficlights

Two trafficlights are simulated.
* trafficlight1: GPIO2-GPIO4 LEDs red, yellow and green
* trafficlight2: GPIO5-GPIO7 LEDs red, yellow and green

## Example

Contents of GPIO.json:

```
[
17, 18,
27, 22, 23,
24, 25, 4
]```

```
17 - GPIO0 (unused)
18 - GPIO1 (unused)
27 - GPIO2 is BCM pin 27 which is trafficlight1 red LED
22 - GPIO3 is BCM pin 22 which is trafficlight1 yellow LED
...
25 - GPIO6 is BCM pin 25 which is trafficlight2 yellow LED
 4 - GPIO7 is BCM pin 4 which is trafficlight2 green LED
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
