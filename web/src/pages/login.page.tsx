import React, { useState } from "react";
import { Button, Container, Paper, TextField, Typography } from "@mui/material";
import { login } from "../api/auth";
import { useMutation } from "react-query";
import { useNavigate } from "react-router-dom";

const LoginPage = () => {
  const navigate = useNavigate();
  const { mutate, isLoading, isError } = useMutation("login", login);

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");

  const handleLogin = () => {
    mutate({
      username,
      password,
    });
  };

  return (
    <Container maxWidth="xs">
      <Paper elevation={3} style={{ padding: "20px", marginTop: "50px" }}>
        <Typography variant="h5" component="h2">
          Login
        </Typography>
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
