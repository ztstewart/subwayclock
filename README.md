# Subwayclock

Subwayclock is a simple program used to power a Raspberry PI-based countdown
clock. The display is similar to the countdown clocks that can be seen on the
New York Subway platform.

![alt text](image.go "Image of the subway clock on a Raspberry Pi Zero)

# Requirements

In order to use this library, you'll need the following:

1. An MTA API key. You can register for an account at: http://datamine.mta.info/
2. A Raspberry PI. I have personally tested this on a Raspberry PI Zero W, but other models should work.
3. An [Inky pHAT](https://shop.pimoroni.com/products/inky-phat) display.

# Usage

1. Modify main.go to use the stop ID for your home station and / or change the labels.
2. Add your MTA API Key in the `client.NewNYCTA(&client.Config{....})` line.
3. Compile the tool: `go build`. You can compile for the Raspberry Pi a more powerful machine: `GOOS=linux GOARCH=arm GOARM=5 go build`
4. (Optional) If you cross compiled (that is, compiled on another machine), copy the binary to your Raspberry Pi.
5. Run the binary: `./subwayclock`. It will continuosly update its information every minute by default.
