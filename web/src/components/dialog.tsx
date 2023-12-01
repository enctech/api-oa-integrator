import * as React from "react";
import { ReactNode } from "react";
import Dialog from "@mui/material/Dialog";
import DialogActions from "@mui/material/DialogActions";
import DialogContent from "@mui/material/DialogContent";
import DialogContentText from "@mui/material/DialogContentText";
import DialogTitle from "@mui/material/DialogTitle";

interface AlertDialogProps {
  isOpen: boolean;
  handleClose: () => void;
  title: string;
  description: string;
  buttons: ReactNode[];
}

export default function AlertDialog(arg: AlertDialogProps) {
  return (
    <>
      <Dialog
        open={arg.isOpen}
        onClose={arg.handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{arg.title}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            {arg.description}
          </DialogContentText>
        </DialogContent>
        <DialogActions>{arg.buttons}</DialogActions>
      </Dialog>
    </>
  );
}
