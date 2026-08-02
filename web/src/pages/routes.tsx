import * as React from "react";
import { styled } from "@mui/material/styles";
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
import HomePage from "./home.page";
import OAConfigsPage from "./config/oa-configs.page";
import OaConfigsDetailsPage from "./config/oa-configs-details.page";
import ThirdPartyConfigsPage from "./config/third-party-configs.page";
import IntegratorConfigsDetailsPage from "./config/third-party-configs-details.page";
import OATransactionPage from "./oa-transactions.page";
import ThirdPartyTransactionsPage from "./third-party-transactions.page";
import { useSession, useSessionGuard } from "../context/session-context";
import AlertDialog from "../components/dialog";
import SettingsIcon from "@mui/icons-material/Settings";
import DashboardIcon from "@mui/icons-material/Dashboard";
import ReceiptLongIcon from "@mui/icons-material/ReceiptLong";
import GroupIcon from "@mui/icons-material/Group";
import LogoutIcon from "@mui/icons-material/Logout";
import LoginIcon from "@mui/icons-material/Login";
import ReceiptIcon from "@mui/icons-material/Receipt";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Button,
} from "@mui/material";
import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import UsersPage from "./users/users.page";
import LogsPage from "./logs.page";
import { AdminOnly } from "../components/auth-guard";
import { useEffect } from "react";

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
  const { session, logout } = useSession();
  const [open, setOpen] = React.useState(session !== null);

  useEffect(() => {
    setOpen(session !== null);
  }, [session]);

  const handleDrawerOpen = () => {
    setOpen(!open);
  };

  const [showLogoutDialog, setShowLogoutDialog] = React.useState(false);

  return (
    <Box sx={{ display: "flex" }}>
      <CssBaseline />
      <AppBar position="fixed" open={open}>
        <Toolbar>
          <AdminOnly>
            <IconButton
              color="inherit"
              aria-label="open drawer"
              edge="start"
              onClick={handleDrawerOpen}
            >
              <MenuIcon style={{ color: "#141617" }} />
            </IconButton>
          </AdminOnly>
          <Typography
            variant="h6"
            noWrap
            sx={{ flexGrow: 1, textTransform: "uppercase" }}
            component="div"
            style={{ color: "#141617" }}
          >
            Online Authorization Dashboard
          </Typography>
          {!session && (
            <Button color="inherit" onClick={() => navigation("/login")}>
              <LoginIcon style={{ color: "#141617" }} />
              <div className="w-2" />
              <Typography style={{ color: "#141617" }}>Login</Typography>
            </Button>
          )}
        </Toolbar>
      </AppBar>
      <Main open={open} className="h-max">
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
                  icon: <DashboardIcon style={{ color: "white" }} />,
                },
                {
                  text: "Configuration",
                  icon: <SettingsIcon style={{ color: "white" }} />,
                  groups: [
                    {
                      text: "Online Authorisation",
                      link: "/oa-configs",
                    },
                    {
                      text: "3rd parties",
                      link: "/3rd-party-config",
                    },
                  ],
                },
                {
                  text: "Transactions",
                  icon: <ReceiptIcon style={{ color: "white" }} />,
                  groups: [
                    {
                      text: "Online Authorisation",
                      link: "/oa-transactions",
                    },
                    {
                      text: "3rd parties",
                      link: "/3rd-party-transactions",
                    },
                  ],
                },
                {
                  text: "Logs",
                  icon: <ReceiptLongIcon style={{ color: "white" }} />,
                  link: "/logs",
                },
                {
                  icon: <GroupIcon style={{ color: "white" }} />,
                  text: "Users",
                  link: "/users",
                },
              ]
            : [
                {
                  text: "Home",
                  link: "/",
                  icon: <DashboardIcon style={{ color: "white" }} />,
                },
                {
                  text: "Transactions",
                  icon: <ReceiptIcon style={{ color: "white" }} />,
                  groups: [
                    {
                      text: "Online Authorisation Transactions",
                      link: "/oa-transactions",
                    },
                    {
                      text: "3rd parties",
                      link: "/3rd-party-transactions",
                    },
                  ],
                },
                {
                  text: "Logs",
                  icon: <ReceiptLongIcon style={{ color: "white" }} />,
                  link: "/logs",
                },
              ]
          ).map(({ text, link, icon, groups }, index) => (
            <ListItem key={`${text}${link}`} disablePadding>
              {groups ? (
                <Accordion
                  className="w-full"
                  elevation={0}
                  sx={{
                    backgroundColor: "#3d4146",
                  }}
                >
                  <AccordionSummary
                    expandIcon={<ExpandMoreIcon style={{ color: "white" }} />}
                    aria-controls="panel1a-content"
                    id="panel1a-header"
                  >
                    {icon && (
                      <>
                        {icon}
                        <div className="w-2" />
                      </>
                    )}
                    <Typography style={{ color: "white" }}>{text}</Typography>
                  </AccordionSummary>
                  <AccordionDetails sx={{ backgroundColor: "#9399a1" }}>
                    {groups.map(({ text, link }, index) => (
                      <ListItemButton
                        key={`${text}${link}`}
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
                  {icon && (
                    <>
                      {icon}
                      <div className="w-2" />
                    </>
                  )}
                  <ListItemText primary={text} style={{ color: "#FFFFFF" }} />
                </ListItemButton>
              )}
            </ListItem>
          ))}
        </List>
        {session && (
          <List
            style={{ position: "absolute", bottom: "0", right: "0", left: "0" }}
          >
            <ListItem disablePadding>
              <ListItemButton
                onClick={() => setShowLogoutDialog(true)}
                sx={{ backgroundColor: "#3d4146" }}
              >
                <LogoutIcon style={{ color: "white" }} />
                <div className="w-2" />
                <ListItemText primary={"Logout"} style={{ color: "#FFFFFF" }} />
              </ListItemButton>
            </ListItem>
          </List>
        )}
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
              navigation("/login");
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
      <Route
        path="/"
        element={
          <AdminGuard>
            <HomePage />
          </AdminGuard>
        }
      />
      <Route path="/login" element={<LoginPage />} />
      {/*<Route path="users/*" element={<Users/>}/>*/}
      <Route
        path="/logs"
        element={
          <AdminGuard>
            <LogsPage />
          </AdminGuard>
        }
      />
      <Route path="/oa-configs" element={<OAConfigsPage />} />
      <Route path="/oa-configs/:id" element={<OaConfigsDetailsPage />} />
      <Route path="/3rd-party-config" element={<ThirdPartyConfigsPage />} />
      <Route
        path="/3rd-party-config/:id"
        element={<IntegratorConfigsDetailsPage />}
      />
      <Route
        path="/oa-transactions"
        element={
          <AdminGuard>
            <OATransactionPage />
          </AdminGuard>
        }
      />
      <Route path="/users" element={<UsersPage />} />
      <Route
        path="/3rd-party-transactions"
        element={
          <AdminGuard>
            <ThirdPartyTransactionsPage />
          </AdminGuard>
        }
      />
    </Routes>
  );
}

const AdminGuard = ({ children }: any) => {
  useSessionGuard();
  return children;
};

export default function TopRoutes() {
  return (
    <BrowserRouter>
      <PersistentDrawerRight />
    </BrowserRouter>
  );
}
