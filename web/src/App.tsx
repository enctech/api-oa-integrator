import React from "react";
import AppRoutes from "./pages/routes";
import { QueryClient, QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from "react-query/devtools";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { SessionProvider } from "./context/session-context";
import { ThemeProvider } from "@emotion/react";
import { createTheme } from "@mui/material";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={theme}>
        <LocalizationProvider dateAdapter={AdapterDayjs}>
          <SessionProvider>
            <ReactQueryDevtools /> {/* Optional devtools */}
            <AppRoutes />
          </SessionProvider>
        </LocalizationProvider>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

const theme = createTheme({
  typography: {
    fontFamily: ["Bai Jamjuree"].join(","),
  },
});

export default App;
