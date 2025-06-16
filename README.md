
# Serial Plotter

![Demo](./documentation/gifs/Screencast.gif)

## Purpose

The arduino serial plotter for the newer versions have had a sever regression including:
 - Allowing only 50 datapoints to be displayed compared to the previous 200
 - Not releasing the serial port resources when the serial plotter is closed
 - Not storing the board type and baud as a user preference

I use the serial plotter when helping write and debug arduino scripts and with the new arduino IDE I spend a lot of time on repetitive steps. I also would like more than 50 data points displayed at a time. This application is to help overcome this downfalls. 

## Features
 - No software defined limit on data points plotted (as many as the hardware can handle).
 - Serial port resource is released when the application is not graphing (allowing other applications to upload)
 - Mobile friendly -- this allows you to view the serial data from an arduino with your phone instead of dragging a laptop around!
 - Filtering functions -- apply causal filters to help with data plotted from noisy sensors


## Development

### Building and Running

`go run .`

### Building an android APK

Download and install android-studio. If you do not want to download the full android studio IDE then just download the [NDK](https://developer.android.com/ndk/downloads). Make sure the `ANDROID_HOME` and `ANDROID_NDK_HOME` environment variables are set to the correct locations. The default install paths are below:

```bash
# Android studio NDK variables
export ANDROID_HOME=$HOME/Android/Sdk
export ANDROID_NDK_HOME=$HOME/Android/Sdk/ndk/29.0.13599879
```

> Note make sure to update the ndk version to the version you have downloaded

Finally run:

```bash
./build-android.sh
```

This will create a _Serial_Plotter.apk_ that can be installed on an android device or emulator.

## Future ideas
 - [ ] serial plotter allows you to save data
 - [ ] serial plotter mobile app
 - [ ] serial plotter + uC project for DIY sensors
 - [ ] allow multiple data inputs
 - [ ] add gauges for one dimensional values
 - [ ] add raw value displays for things measurements like temperature
 - [ ] add more filters


## Items that need done but aren't quite bugs
 - [x] rename pseudo source to dummy source
 - [x] rename pseudo.Transform to pseudo.Function
 - [x] Add clear button to reset data
 - [ ] Have selected options persist when re-opening application (using fyne preferences)
 - [ ] Create custom icon instead of using arduinos icon


## Bugs to fix
 - [ ] Fix Y axis display when the data range is zero
 - [x] Don't let the user press start twice
 - [ ] Figure out how to use the theme to set the graph colors...
 - [x] Gracefully handle serial port errors if port can't be opened

