import { Navigate, Route, Routes } from "react-router-dom";
import "./App.css";
import Login from "./pages/Login";
import Register from "./pages/Register";
import NotFound from "./pages/Notfound.tsx";
import Dashboard from "./pages/Dashboard";
import { useState } from "react";
import { LOGIN_LS_KEY } from "./constants";

function App() {
  const [loggedIn, setLoggedIn] = useState<boolean>(
    localStorage.getItem(LOGIN_LS_KEY) !== null
  );

  return (
    <div className="App">
      <Routes>
        <Route path="/" element={<Navigate to="/login" />} />
        <Route
          path="/login"
          element={
            loggedIn ? (
              <Navigate to="/dashboard" />
            ) : (
              <Login setLoggedIn={setLoggedIn} />
            )
          }
        />
        <Route
          path="/register"
          element={loggedIn ? <Navigate to="/dashboard" /> : <Register />}
        />
        <Route
          path="/dashboard"
          element={
            loggedIn ? (
              <Dashboard setLoggedIn={setLoggedIn} />
            ) : (
              <Navigate to="/login" />
            )
          }
        />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </div>
  );
}

export default App;
