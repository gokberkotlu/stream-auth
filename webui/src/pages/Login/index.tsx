import { Button, message } from "antd";
import React, { useEffect, useRef, useState } from "react";
import { Link, useNavigate } from "react-router-dom";

interface Props {
  setLoggedIn: React.Dispatch<React.SetStateAction<boolean>>;
}

interface IRecRes {
  name: string;
  rect: {
    Min: {
      X: number;
      Y: number;
    };
    Max: {
      X: number;
      Y: number;
    };
  };
}

const Login: React.FC<Props> = ({ setLoggedIn }) => {
  const socket = new WebSocket("ws://localhost:8080/face-rec");
  const navigate = useNavigate();

  const [rectanglePoints, setRectanglePoints] = useState<string>();
  const [userName, setUserName] = useState<string>("");

  const videoElement = useRef<HTMLVideoElement>(null);

  const wsMessageInParseError = useRef<string>("");

  const startStream = async () => {
    try {
      // Access the user's media devices
      let localStream = await navigator.mediaDevices.getUserMedia({
        video: true,
        audio: false,
      });

      if (videoElement.current) {
        videoElement.current.srcObject = localStream;
      }

      // WebSocket connection opened
      socket.onopen = () => {
        console.log("WebSocket connection opened");
      };

      // WebSocket connection closed
      socket.onclose = () => {
        console.log("WebSocket connection closed");
      };

      socket.onmessage = (socketMessage) => {
        captureAndSendImage();
        console.log(socketMessage.data, typeof socketMessage.data);
        if (socketMessage?.data) {
          let response: IRecRes;

          try {
            response = JSON.parse(socketMessage.data);
          } catch (err) {
            if (wsMessageInParseError.current !== socketMessage.data) {
              wsMessageInParseError.current = socketMessage.data;
              message.error(socketMessage.data);
            }

            setRectanglePoints("0,0 0,0 0,0 0,0");
            return;
          }

          const { name, rect } = response;

          const { X: x1, Y: y1 } = rect?.Min;
          const { X: x2, Y: y2 } = rect?.Max;

          if (name) {
            setUserName(name);
          }

          if (x1 && y1 && x2 && y2) {
            setRectanglePoints(
              `${x1},${y1} ${x2},${y1} ${x2},${y2} ${x1},${y2}`
            );

            login(name);
          } else {
            setRectanglePoints("0,0 0,0 0,0 0,0");
          }
        } else {
          setRectanglePoints("0,0 0,0 0,0 0,0");
        }
      };

      // WebSocket connection error
      socket.onerror = (error) => {
        console.error("WebSocket error:", error);
      };

      // Capture and send image periodically
      // setInterval(captureAndSendImage, 1000); // Adjust the interval as needed

      videoElement.current?.addEventListener("canplay", () => {
        captureAndSendImage();
      });
    } catch (error) {
      console.error("Error accessing media devices:", error);
    }
  };

  const closeWsConnection = () => {
    // Check if the connection is still open before attempting to close it
    if (socket && socket.readyState === WebSocket.OPEN) {
      socket.close();
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

      // Send the image data to the server via WebSocket
      socket.send(imageData);

      const date = new Date();
      console.log(
        `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()}`
      );
    } catch (err) {
      console.log("captureAndSendImage error:", err);
    }
  };

  const login = (userName: string) => {
    localStorage.setItem("stream_auth_user", userName);
    setLoggedIn(true);
    closeWsConnection();
    navigate("/dashboard");
  };

  useEffect(() => {
    startStream();

    return () => {
      closeWsConnection();
    };
  }, []);

  return (
    <div>
      <h1>Periodic Image Capture</h1>

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
        <span
          style={{
            position: "absolute",
            width: "100%",
            height: "100%",
            left: 0,
            top: 0,
          }}
        >
          {userName}
        </span>

        <br />
        <br />
        <br />
        <Link to="/register">
          <Button type="primary">Register</Button>
        </Link>
      </div>
    </div>
  );
};

export default Login;
