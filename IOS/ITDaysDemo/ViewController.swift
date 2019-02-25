//
//  ViewController.swift
//  ITDaysDemo
//
//  Created by Sorin Sebastian Mircea on 07/11/2018.
//  Copyright © 2018 Sorin Sebastian Mircea. All rights reserved.
//

import UIKit
import SceneKit
import ARKit

import RxCocoa
import RxSwift
import RxGesture
import Alamofire
import Vision

enum Selected {
    case nothing, robot, lamp;
}
class ViewController: UIViewController, ARSCNViewDelegate {
    @IBOutlet var sceneView: ARSCNView!
    @IBOutlet weak var interractButton: UIButton!
    @IBOutlet weak var selectedLabel: UILabel!
    
    private let disposeBag = DisposeBag()
    private var selected = Selected.nothing
    
    var imageHighlightAction: SCNAction {
        return .sequence([
            .wait(duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOpacity(to: 0.96, duration: 0.25),
            .fadeOpacity(to: 0.15, duration: 0.25),
            .fadeOut(duration: 0.5),
            .removeFromParentNode()
            ])
    }
    
    
    override func viewDidLoad() {
        super.viewDidLoad()
        
        // Set the view's delegate
        sceneView.delegate = self

        // Show statistics such as fps and timing information
        sceneView.showsStatistics = false
        sceneView.debugOptions = [ARSCNDebugOptions.showFeaturePoints]

        // Create a new scene
        let scene = SCNScene(named: "art.scnassets/ship.scn")!

        // Set the scene to the view
        sceneView.scene = scene

        // Set up button interaction
        interractButtonAction()
    }
    
    override func viewWillAppear(_ animated: Bool) {
        super.viewWillAppear(animated)
        
        setARconfig()
    }
    
    func setARconfig() {
        // Create a session configuration
        let configuration = ARWorldTrackingConfiguration()
        
        // Object detection
        configuration.detectionObjects = ARReferenceObject.referenceObjects(inGroupNamed: "AR Object Detection", bundle: Bundle.main)!
        
        // Image detection
        configuration.detectionImages = ARReferenceImage.referenceImages(inGroupNamed: "AR Image Detection", bundle: Bundle.main)!
        
        // Run the view's session
        sceneView.session.run(configuration, options: [.resetTracking, .removeExistingAnchors])
    }
    
    // Object recognition
    func renderer(_ renderer: SCNSceneRenderer, nodeFor anchor: ARAnchor) -> SCNNode? {
        let node = SCNNode()
        
        if let objectAnchor = anchor as? ARObjectAnchor {
            print("Detected!")
            print(objectAnchor.name!)
            
            if(objectAnchor.name! == "robot" || objectAnchor.name! == "robotLeg") {
                self.selected = Selected.robot
            }
            
            if(objectAnchor.name! == "lamp") {
                self.selected = Selected.lamp
            }
            
            let boxGeometry = SCNBox(width: 0.04, height: 0.04, length: 0.04, chamferRadius: 0.0)
            let boxNode = SCNNode(geometry: boxGeometry)
            boxNode.position = SCNVector3( objectAnchor.referenceObject.center.x,
                                           objectAnchor.referenceObject.center.y,
                                           objectAnchor.referenceObject.center.z
            )
            node.addChildNode(boxNode)
        
        }
        
        return node
    }
    
    // Image recognition
    func renderer(_ renderer: SCNSceneRenderer, didAdd node: SCNNode, for anchor: ARAnchor) {
        guard let imageAnchor = anchor as? ARImageAnchor else { return }
        let referenceImage = imageAnchor.referenceImage
        let imageName = referenceImage.name ?? "no name"
        
        print("Image found: ", imageName)
        if(imageName == "green") {
            self.selected = Selected.robot
            self.selectedLabel.text = "Robot"
            self.selectedLabel.isHidden = false;
        }
        if(imageName == "red") {
            self.selected = Selected.lamp
            self.selectedLabel.text = "Lamp"
             self.selectedLabel.isHidden = false;
        }
        
        // Create a plane to visualize the initial position of the detected image.
        let plane = SCNPlane(width: referenceImage.physicalSize.width,height: referenceImage.physicalSize.height)
        let planeNode = SCNNode(geometry: plane)
        planeNode.opacity = 0.25
        planeNode.eulerAngles.x = -.pi / 2
        planeNode.runAction(self.imageHighlightAction)
        node.addChildNode(planeNode)
        
        // Reset the tracking session
        DispatchQueue.main.asyncAfter(deadline: .now() + 5) {
            print("Removed anchor")
             self.selectedLabel.isHidden = true;
            self.setARconfig()
        }
        
    }
    
    
    override func viewWillDisappear(_ animated: Bool) {
        super.viewWillDisappear(animated)
        
        // Pause the view's session
        sceneView.session.pause()
    }

    // MARK: - ARSCNViewDelegate
    
/*
    // Override to create and configure nodes for anchors added to the view's session.
    func renderer(_ renderer: SCNSceneRenderer, nodeFor anchor: ARAnchor) -> SCNNode? {
        let node = SCNNode()
     
        return node
    }
*/
    
    func session(_ session: ARSession, didFailWithError error: Error) {
        // Present an error message to the user
        
    }
    
    func sessionWasInterrupted(_ session: ARSession) {
        // Inform the user that the session has been interrupted, for example, by presenting an overlay
        
    }
    
    func sessionInterruptionEnded(_ session: ARSession) {
        // Reset tracking and/or remove existing anchors if consistent tracking is required
        
    }
    
    func interractButtonAction() {
        // Tap
        interractButton.rx
            .tapGesture(configuration: { gestureRecognizer, delegate in
                delegate.simultaneousRecognitionPolicy = .never
            })
            .subscribe({ev in
                print("Sending request to robot")
                let url = URL(string: "http://192.168.43.24:8000/writeITDAYS")
                Alamofire.request(url!)
            })
            .disposed(by: disposeBag)
        
        // Long Press Gesture
        interractButton.rx
            .longPressGesture(configuration: { gestureRecognizer, delegate in
                delegate.simultaneousRecognitionPolicy = .never
            })
            .when(.began)
            .subscribe(onNext: {ev in
                print("Sending request to lamp")
                let headers: HTTPHeaders = [
                    "Authorization": "##",
                    "Accept": "application/json"
                ]

                Alamofire.request("https://api.lifx.com/v1/lights/all/state?power=on&brightness=1.0&fast=true",  method: .put, headers: headers)
            })
            .disposed(by: disposeBag)
    }
}
