import { Button, Input, message } from "antd";
import React, { useEffect, useRef, useState } from "react";
import { Link } from "react-router-dom";
import { HomeOutlined } from "@ant-design/icons";

interface ISingleRecRes {
  Min: {
    X: number;
    Y: number;
  };
  Max: {
    X: number;
    Y: number;
  };
}

// const socketRef.current = new WebSocket("ws://localhost:8080/face-register?name=thename");

const Register: React.FC = () => {
  const [rectanglePoints, setRectanglePoints] = useState<string>();
  const [userName, setUserName] = useState<string>("");

  const socketRef = useRef<WebSocket | null>(null);
  const registerCount = useRef<number>(5);

  const videoElement = useRef<HTMLVideoElement>(null);

  const startStream = async () => {
    const regex = /^[^\s]*$/;
    if (!regex.test(userName)) {
      message.error("Username cannot contain spaces");
      return;
    } else if (userName.length < 3) {
      message.error("Username must be at least 3 characters");
      return;
    }

    try {
      socketRef.current = new WebSocket(
        `ws://localhost:8080/face-register?name=${userName}`
      );

      // Access the user's media devices
      let localStream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: false,
      });

      if (videoElement.current) {
        videoElement.current.srcObject = localStream;
      }

      // WebSocket connection opened
      socketRef.current.onopen = () => {
        console.log("WebSocket connection opened");
      };

      // WebSocket connection closed
      socketRef.current.onclose = () => {
        console.log("WebSocket connection closed");
      };

      socketRef.current.onmessage = (socketMessage) => {
        if (registerCount.current > 1) {
          captureAndSendImage();
          console.log(socketMessage.data, typeof socketMessage.data);
          if (socketMessage?.data) {
            let rect: ISingleRecRes;

            try {
              rect = JSON.parse(socketMessage.data);
            } catch (err) {
              message.error(socketMessage.data);
              setRectanglePoints("0,0 0,0 0,0 0,0");
              return;
            }

            const { X: x1, Y: y1 } = rect?.Min;
            const { X: x2, Y: y2 } = rect?.Max;

            if (x1 && y1 && x2 && y2) {
              setRectanglePoints(
                `${x1},${y1} ${x2},${y1} ${x2},${y2} ${x1},${y2}`
              );
              registerCount.current -= 1;
            } else {
              setRectanglePoints("0,0 0,0 0,0 0,0");
            }
          } else {
            setRectanglePoints("0,0 0,0 0,0 0,0");
          }
        } else {
          setRectanglePoints("0,0 0,0 0,0 0,0");
        }
      };

      // WebSocket connection error
      socketRef.current.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      videoElement.current?.addEventListener("canplay", () => {
        captureAndSendImage();
      });
    } catch (error) {
      console.error("Error accessing media devices:", error);
    }
  };

  const closeWsConnection = () => {
    // Check if the connection is still open before attempting to close it
    if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
      socketRef.current.close();
    }
  };

  // Function to capture and send image data
  const captureAndSendImage = () => {
    try {
      console.log("send image");
      const canvas = document.createElement("canvas");
      const context = canvas.getContext("2d");

      if (videoElement.current) {
        canvas.width = videoElement.current.videoWidth;
        canvas.height = videoElement.current.videoHeight;

        context?.drawImage(
          videoElement.current,
          0,
          0,
          canvas.width,
          canvas.height
        );
      }

      const imageData = canvas.toDataURL("image/jpeg");

      if (socketRef.current) {
        // Send the image data to the server via WebSocket
        socketRef.current.send(imageData);
      }

      const date = new Date();
      console.log(
        `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
      );
    } catch (err) {
      console.log("captureAndSendImage error:", err);
    }
  };

  const startRegistration = () => {
    startStream();
  };

  const stopRegistration = () => {
    closeWsConnection();
  };

  useEffect(() => {
    return () => closeWsConnection();
  }, []);

  return (
    <div style={{ height: "100vh", backgroundColor: "#c3c3c3" }}>
      <div
        className="videoRect"
        style={{ position: "relative", display: "inline-block" }}
      >
        <video ref={videoElement} autoPlay></video>

        <svg
          style={{
            position: "absolute",
            width: "100%",
            height: "100%",
            left: 0,
            top: 0,
          }}
        >
          <polygon
            points={rectanglePoints}
            style={{
              fill: "transparent",
              stroke: "#DC143C",
              strokeWidth: 2,
            }}
          />
        </svg>
      </div>
      <br />
      <br />
      <div
        style={{
          display: "inline-block",
        }}
      >
        <Input
          style={{ width: 400, height: 60 }}
          placeholder="Enter username"
          value={userName}
          onChange={(e) => {
            setUserName(e.target.value);
          }}
        />
        <div
          style={{
            display: "flex",
            justifyContent: "center",
            alignItems: "center",
            marginTop: 15,
          }}
        >
          <Button
            type="primary"
            onClick={startRegistration}
            danger
            style={{ marginLeft: 5, marginRight: 5, width: "100%", height: 45 }}
          >
            Start Registration
          </Button>
          <Button
            type="primary"
            onClick={stopRegistration}
            danger
            style={{ marginRight: 5, width: "100%", height: 45 }}
          >
            Stop Registration
          </Button>
        </div>
      </div>

      <Link to="/">
        <Button
          type="primary"
          shape="circle"
          size="large"
          style={{
            position: "absolute",
            top: 20,
            left: 30,
          }}
        >
          <HomeOutlined />
        </Button>
      </Link>
    </div>
  );
};

export default Register;
