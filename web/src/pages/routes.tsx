import * as React from "react";
import { styled, useTheme } from "@mui/material/styles";
import Box from "@mui/material/Box";
import Drawer from "@mui/material/Drawer";
import LoginPage from "./login.page";
import MuiAppBar, { AppBarProps as MuiAppBarProps } from "@mui/material/AppBar";
import Toolbar from "@mui/material/Toolbar";
import CssBaseline from "@mui/material/CssBaseline";
import List from "@mui/material/List";
import Typography from "@mui/material/Typography";
import IconButton from "@mui/material/IconButton";
import MenuIcon from "@mui/icons-material/Menu";
import ListItem from "@mui/material/ListItem";
import ListItemButton from "@mui/material/ListItemButton";
import ListItemText from "@mui/material/ListItemText";
import { BrowserRouter, Route, Routes, useNavigate } from "react-router-dom";
import HomePage from "./home/home.page";
import LogsPage from "./logs.page";
import OAConfigsPage from "./config/oa-configs.page";
import OaConfigsDetailsPage from "./config/oa-configs-details.page";
import IntegratorConfigsPage from "./config/integrator-configs.page";
import IntegratorConfigsDetailsPage from "./config/integrator-configs-details.page";
import OATransactionPage from "./oa-transactions.page";
import IntegratorTransactionsPage from "./integrators-transactions.page";
import { useSession } from "../context/session-context";
import AlertDialog from "../components/dialog";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Button,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import UsersPage from "./users.page";

const drawerWidth = 240;

const Main = styled("main", { shouldForwardProp: (prop) => prop !== "open" })<{
  open?: boolean;
}>(({ theme, open }) => ({
  flexGrow: 1,
  padding: theme.spacing(3),
  transition: theme.transitions.create("margin", {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  marginLeft: 0,
  ...(open && {
    transition: theme.transitions.create("margin", {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: drawerWidth,
  }),
  position: "relative",
}));

interface AppBarProps extends MuiAppBarProps {
  open?: boolean;
}

const AppBar = styled(MuiAppBar, {
  shouldForwardProp: (prop) => prop !== "open",
})<AppBarProps>(({ theme, open }) => ({
  transition: theme.transitions.create(["margin", "width"], {
    easing: theme.transitions.easing.sharp,
    duration: theme.transitions.duration.leavingScreen,
  }),
  backgroundColor: "#fdc300",
  ...(open && {
    width: `calc(100% - ${drawerWidth}px)`,
    transition: theme.transitions.create(["margin", "width"], {
      easing: theme.transitions.easing.easeOut,
      duration: theme.transitions.duration.enteringScreen,
    }),
    marginLeft: drawerWidth,
  }),
}));

const DrawerHeader = styled("div")(({ theme }) => ({
  display: "flex",
  alignItems: "center",
  padding: theme.spacing(0, 1),
  // necessary for content to be below app bar
  ...theme.mixins.toolbar,
  justifyContent: "flex-end",
}));

function PersistentDrawerRight() {
  const navigation = useNavigate();
  const theme = useTheme();
  const [open, setOpen] = React.useState(true);

  const handleDrawerOpen = () => {
    setOpen(!open);
  };

  const { session, logout } = useSession();

  const [showLogoutDialog, setShowLogoutDialog] = React.useState(false);

  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <AppBar position="fixed" open={open}>
        <Toolbar>
          <IconButton
            color="inherit"
            aria-label="open drawer"
            edge="start"
            onClick={handleDrawerOpen}
          >
            <MenuIcon style={{ color: "#141617" }} />
          </IconButton>
          <Typography
            variant="h6"
            noWrap
            sx={{ flexGrow: 1, textTransform: "uppercase" }}
            component="div"
            style={{ color: "#141617" }}
          >
            Online Authorization Dashboard
          </Typography>
          {session ? (
            <Button color="inherit" onClick={() => setShowLogoutDialog(true)}>
              <Typography style={{ color: "#141617" }}>Logout</Typography>
            </Button>
          ) : (
            <Button color="inherit" onClick={() => navigation("/login")}>
              <Typography style={{ color: "#141617" }}>Login</Typography>
            </Button>
          )}
        </Toolbar>
      </AppBar>
      <Main open={open}>
        <DrawerHeader />
        <AppRoutes />
      </Main>
      <Drawer
        PaperProps={{
          sx: {
            backgroundColor: "#3d4146",
          },
        }}
        sx={{
          flexShrink: 0,
          "& .MuiDrawer-paper": {
            width: drawerWidth,
          },
        }}
        variant="persistent"
        anchor="left"
        open={open}
      >
        <DrawerHeader />
        <List>
          {(session
            ? [
                {
                  text: "Home",
                  link: "/",
                },
                {
                  text: "Configuration",
                  groups: [
                    {
                      text: "Online Authorisation",
                      link: "/oa-configs",
                    },
                    {
                      text: "Integrators",
                      link: "/integrator-configs",
                    },
                  ],
                },
                {
                  text: "Transactions",
                  groups: [
                    {
                      text: "Online Authorisation",
                      link: "/oa-transactions",
                    },
                    {
                      text: "Integrators",
                      link: "/integrator-transactions",
                    },
                  ],
                },
                {
                  text: "Logs",
                  link: "/logs",
                },
                {
                  text: "Users",
                  link: "/users",
                },
              ]
            : [
                {
                  text: "Home",
                  link: "/",
                },
                {
                  text: "Transactions",
                  groups: [
                    {
                      text: "Online Authorisation Transactions",
                      link: "/oa-transactions",
                    },
                    {
                      text: "Integrator Transactions",
                      link: "/integrator-transactions",
                    },
                  ],
                },
                {
                  text: "Logs",
                  link: "/logs",
                },
              ]
          ).map(({ text, link, groups }, index) => (
            <ListItem key={link} disablePadding>
              {groups ? (
                <Accordion
                  className="w-full"
                  elevation={0}
                  sx={{
                    backgroundColor: "#3d4146",
                  }}
                >
                  <AccordionSummary
                    expandIcon={<ExpandMoreIcon />}
                    aria-controls="panel1a-content"
                    id="panel1a-header"
                  >
                    <Typography style={{ color: "#FFFFFF" }}>{text}</Typography>
                  </AccordionSummary>
                  <AccordionDetails sx={{ backgroundColor: "#9399a1" }}>
                    {groups.map(({ text, link }, index) => (
                      <ListItemButton
                        onClick={() => navigation(link)}
                        color="white"
                      >
                        <ListItemText
                          primary={text}
                          style={{ color: "#FFFFFF" }}
                        />
                      </ListItemButton>
                    ))}
                  </AccordionDetails>
                </Accordion>
              ) : (
                <ListItemButton
                  onClick={() => navigation(link)}
                  sx={{ backgroundColor: "#3d4146" }}
                >
                  <ListItemText primary={text} style={{ color: "#FFFFFF" }} />
                </ListItemButton>
              )}
            </ListItem>
          ))}
        </List>
      </Drawer>
      <AlertDialog
        isOpen={showLogoutDialog}
        handleClose={() => setShowLogoutDialog(false)}
        title={"Are you sure you want to logout?"}
        description={"You will be logged out of the system."}
        buttons={[
          <Button
            key="cancel"
            onClick={() => setShowLogoutDialog(false)}
            color="primary"
          >
            Cancel
          </Button>,
          <Button
            key="logout"
            onClick={() => {
              logout();
              navigation("/");
              setShowLogoutDialog(false);
            }}
            color="primary"
            autoFocus
          >
            Logout
          </Button>,
        ]}
      />
    </Box>
  );
}

function AppRoutes() {
  return (
    <Routes>
      <Route path="/" element={<HomePage />} />
      <Route path="/login" element={<LoginPage />} />
      {/*<Route path="users/*" element={<Users/>}/>*/}
      <Route path="/logs" element={<LogsPage />} />
      <Route path="/oa-configs" element={<OAConfigsPage />} />
      <Route path="/oa-configs/:id" element={<OaConfigsDetailsPage />} />
      <Route path="/integrator-configs" element={<IntegratorConfigsPage />} />
      <Route
        path="/integrator-configs/:id"
        element={<IntegratorConfigsDetailsPage />}
      />
      <Route path="/oa-transactions" element={<OATransactionPage />} />
      <Route path="/users" element={<UsersPage />} />
      <Route
        path="/integrator-transactions"
        element={<IntegratorTransactionsPage />}
      />
    </Routes>
  );
}

export default function TopRoutes() {
  return (
    <BrowserRouter>
      <PersistentDrawerRight />
    </BrowserRouter>
  );
}
