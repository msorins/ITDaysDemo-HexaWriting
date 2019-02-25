# ItDaysDemo IOS + ROBOT

## Youtube demo video
[![](http://img.youtube.com/vi/5WrTGGZzSMA/0.jpg)](http://www.youtube.com/watch?v=5WrTGGZzSMA "Youtube demo video")

## Screens
![Demo Image ](https://github.com/msorins/ITDaysDemo-HexaWriting/blob/master/1.png?raw=true "Demo Image")

![Demo Image ](https://github.com/msorins/ITDaysDemo-HexaWriting/blob/master/2.png?raw=true "Demo Image")

# Idea
*Connected devices*, *IOT* are the main concepts behind this demo.

Using the camera from an IOS device, the user can point either to a smart lamp, or to the hexa robot to invoke an action.

The smart lamp would simply turn on, but the robot is going to start writing *IT DAYS*

# IOS
I have used ARCore for image and object detection.


These manually drawn patterns are used for image detection:
* ![Demo Image ](https://github.com/msorins/ITDaysDemo-HexaWriting/blob/master/3.png?raw=true "Demo Image")
* ![Demo Image ](https://github.com/msorins/ITDaysDemo-HexaWriting/blob/master/4.jpg?raw=true "Demo Image")

# ROBOT
In order to control it, an API server is running on the robot, listening to a certain trigger (that comes from the tablet) to start writing a sentence. As an addition, there are also written a bunch of endpoints for manually controlling the robot (from a web interface)

The animations were made with a helper program called *Hexa Simulator*, and then imported as strings into the program

> To do: add picture


# Technologies used

* GO


> The Project was realised for the conference [ITDAYS](https://www.itdays.ro/speaker/sorin-sebastian).