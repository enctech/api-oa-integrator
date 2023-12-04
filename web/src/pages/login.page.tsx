import React, { useEffect, useState } from "react";
import { Button, Container, Paper, TextField } from "@mui/material";
import { login } from "../api/auth";
import { useMutation } from "react-query";
import { useNavigate } from "react-router-dom";
import { useSession } from "../context/session-context";

const LoginPage = () => {
  const navigate = useNavigate();
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const { mutate, data } = useMutation(["login", username, password], login);
  const { login: loginUser } = useSession();

  const handleLogin = () => {
    mutate({
      username,
      password,
    });
  };

  useEffect(() => {
    if (!data) return;
    loginUser(data);
    navigate("/");
  }, [data]);

  return (
    <Container maxWidth="xs">
      <Paper
        elevation={3}
        style={{
          padding: "20px",
          marginTop: "50px",
        }}
        className="items-center justify-center"
      >
        <div className="flex">
          <div className="flex-grow" />
          <img
            src={"images/logo_enctech.svg"}
            alt="Description"
            width="100"
            height="100"
          />
          <div className="flex-grow" />
        </div>
        <form>
          <TextField
            fullWidth
            label="Username"
            variant="outlined"
            margin="normal"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
          <TextField
            fullWidth
            label="Password"
            type="password"
            variant="outlined"
            margin="normal"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
          <div className="h-4" />
          <Button
            variant="contained"
            color="primary"
            fullWidth
            onClick={handleLogin}
          >
            Login
          </Button>
          <div className="h-4" />
          <Button
            variant="contained"
            sx={{ backgroundColor: "darkgray", color: "#fff" }}
            fullWidth
            onClick={() => navigate("/")}
          >
            Back to Home
          </Button>
        </form>
      </Paper>
    </Container>
  );
};

export default LoginPage;
