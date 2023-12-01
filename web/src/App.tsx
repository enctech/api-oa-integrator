import React from "react";
import AppRoutes from "./pages/routes";
import { QueryClient, QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from "react-query/devtools";
import { LocalizationProvider } from "@mui/x-date-pickers";
import { AdapterDayjs } from "@mui/x-date-pickers/AdapterDayjs";
import { SessionProvider } from "./context/session-context";

const queryClient = new QueryClient();

function App() {
  return (
    <QueryClientProvider client={queryClient}>
      <LocalizationProvider dateAdapter={AdapterDayjs}>
        <SessionProvider>
          <ReactQueryDevtools /> {/* Optional devtools */}
          <AppRoutes />
        </SessionProvider>
      </LocalizationProvider>
    </QueryClientProvider>
  );
}

export default App;
