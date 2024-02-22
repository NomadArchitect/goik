# Roadmap

TODO:

1. Octapod (two longer limbs that point backwards)
1. Specify segment lengths  and run IK for finding angles that make sure that all
   legs touch the ground. Update rest angles from this and update presets (and new command "ground" to update new rest angles)
1. Specify leg num when modifying individual legs
1. Broadcast new positions (sync write?)

## R & D

1. Add new solver. SolveBodyIK
1. Add a new transformer for the body frame. Rotate X/Y/Z. Translate X/Y/Z. This maps nicely to pitch/yaw/roll and shifting the center of gravity 
1. Translate, solve and define these as the new rest angles ! and we can use this to add "funny walks"
1. Add an active body definition (that can be modified) and a reset body definition

1. SolveBodyIK needs the following parameters:
 end effector coordinate + rotate XYZ + translate XYZ

## GOIK

1. X-march / arc pattern (walk in a circle. Not rotation). Piece de resistance: Arc with spin ?
1. Define an intermediat sequence (_add_ a new and faster revert (also with two phases)) to use between gaits -> a) bring all legs down. b) Lift one leg at a time, rotate coxa to neutral and then bring femur and tibia to neutral angles
1. Reading / saving body definition. Add a list of coxa angles. This will allow for uneven distribution of legs around the robot body.
1. Remove hard coded body definitions and instead make a folder with predefined bodies (drop the reset command)
1. Add pitch/yaw/roll of robot body. Ref notes (rotate coxas in X/Y/Z without relocating effectors)
1. Script / interpreter functionality (With regard to scripting. Let update send a message to a channel to signal that the next command can execute)
1. Add the number of servos and size (TLV style) in exported files.
1. Document file format
1. Parametrize number of interpolation steps (have to be odd)
1. cmd: listing all saved pods
1. Motionplanner (Stretch)
1. Get list of stored moves / primitives from controller (via UDP)

### Refactor
1. Refactor members that should be private
1. Ismoving etc => states

## BUGS
- When stopping mid air after recording a primitive. Is this really a bug or should we just use a custom/shorter arc ? (separate start / mid / end cycle?)
1. debug channel not set after gait change => hang
1. Check for valid solution (regarding stretching og segments) for up/down/pitch/yaw/roll etc
1. 	- revert hensyntar ikke endringer fra "ground"
1. - validering på "ground"



## Documentation
1. Trigger motion primitive + number of repetitions
1. Example protocol
1. State chart of how to use primitives
1. Consider hepta/penta as experimental. Do a dynamic check for the number of legs and gait pattern. Number of indices in the gait-pattern will depend on the number of legs

## Axioms
1. All transitions will be ok if the end state of one equals the start state of another
1. one cycle's end condition can be next cycle's start condition. Add "repeat" command.
1. Primitives can be chained together seamlessly (maybe with the help of intermediates)


## Servo configuration / ESP
1.  Test the effect of punch-parameter in the Dynamixel wizard
1. PID tuning
1. Set delayed params. Use broadcast
1. Compile with correct log level (to remove logging in driver)
1. Refactor out all "example" code
1. Spiffs image with primitives (test)


## Building
1. Model and test the big hexapod
1. Posters + video for maker faire
1. Nicer shell
1. Make STEP-files available
1. Make small hexapod with PWM servos as a demo


## Next video
- Tittel: NO CODE Hexapod / MAJOR UPDATE (no longer need to know IK/FK)
- Si noe om udda mønster. Monterte coxa's feil vei
- Side by side video av bevegelser ?
- Dragon skin tips, reversing, recording, normalizing, custom coxa separation, primitives (and how to use them)
- shots: wave, ripple, tripod, rotate, pitch, yaw, roll, speed change, custom moves, walking in N/S/E/W/NE/NW/NE/SW/SE, arc
- No longer any need to struggle with kinematics (forward or inverse) to build and run a hexapod. No need to understand kinematic diagrams or denavit hartenberg parameters
- no longer any need to understand how to implement gaits,
- servo type agonsitic. Only need to translate a number 0-1024 (512 == 0 degrees) to either PWM or use as raw data.
- separate video showing design steps. Use big hexapod as case study
  Husk thumbnail-tool (https://new.express.adobe.com/)
- video shots:
  - show example 0 with tripod and wave gait
  - show & tell: only anchor points and rest angles. Demonstrate with eight legged bot with two longer legs pointing backwards
  - spiffs + motion primitives. Now running in loop, but could equally well be triggered by remote, or over the network
  - show octapod example
  - forskjellige vinkler / segmenter på forskjellige bein. Endre en enkelt eller flere. IK funker fremdeles
  - ground-kommando
  - spider
  - si noe om størrelsen på recording-primitivene (avh av interpolasjonsstep)

  


  ### Video demo 4

1) Eksempel på hvordan lese primitiv
2) Si noe om endianness, chaining, start-endstate etc ved chaining => mer komplekse ting
	spiffskode

3) Lag eksempel på chaining og kjør gjennom
	
4) Why did it walk funny in the last video ?
	My coxas were upside down


  ### Vide demo 5
            - "But, In your previous video, your pod looked like it had parkinson and also had suffered a stroke", you might say
            - Well. It turns out I had mounted my coxa servos upside down (hence the new bitmask option in the export command) and flooding
             the microcontroller with a virtual packet storm wasn't very condusive to fluid movement either
             Vis video av pakkestorm
             - Using primitives, the motion is way more fluid, but there is still some way to go
              
TODO
        - So, what's next on the TODO-list ?
        - Im now sending individual "Goal Position" commands to each servo in sequence. That means that there is a small time delay
        - between when the first and last servos receive ther commands to move.
        - I'll modify the primitive parsing and dynamixel control code to instead broadcast a single packet on the dynamxiel
          networks for each "frame", so that all servos recevie their instruction at the same time
        - There is also a difference between operating a servo stand alone and operating it in an assembly. By default, the 
          XL-320 servos only have the Proportional component set as non zero. There is some room for PID tuning. (The Coxa servos may
          require a different set of PID parameters compared to the femur and tibia parameters for example)
        - All inverse kinematics so far have focused on solving the equations for having the end effectors where I want them to be
        - Next up will be to have the end effectors anchored and then solve the IK equations for the coxa origins. This will allow me
          to implement pitch/yaw and roll in addition to translating the body in XYZ in reference to the end effector positions
        - So far I have only played around with translations a bit. This allows me to alter the "ride height" when it moves. 
        - After I have lined up my homogeneous transformation matrices properly it should also be possible to combine this with gaits, so that 
          pitch/yaw/roll/translation moves can be combined with gaits. No promises, but I'll give it a try :)


        - 






## Commercial
- Free version. Not able to save or export
- Paid version (distribution ? How to lock to computer ?)
- Rebrand to  hexapod designer ?
- If open source: GPL3

