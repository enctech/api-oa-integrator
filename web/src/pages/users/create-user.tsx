import { Button, Checkbox, FormControlLabel, TextField } from "@mui/material";
import React from "react";
import DialogTitle from "@mui/material/DialogTitle";
import Dialog from "@mui/material/Dialog";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import { SubmitHandler, useForm } from "react-hook-form";
import { useMutation } from "react-query";
import { createUser } from "../../api/auth";
import DialogContentText from "@mui/material/DialogContentText";

interface CreateUserFormData {
  name: string;
  username: string;
  password: string;
  isAdmin: boolean;
}

interface Props {
  isVisible: boolean;
  close: () => void;
  onShowError: (message: string) => void;
  onUserCreated: () => void;
}

export default function CreateUser(arg: Props) {
  const { register, handleSubmit, watch, setValue } =
    useForm<CreateUserFormData>({
      defaultValues: {
        name: "",
        password: "",
        isAdmin: false,
        username: "",
      },
    });

  const { mutate, data: newData } = useMutation("createUser", createUser, {
    onSettled: () => {
      arg.onUserCreated();
      arg.close();
    },
    onError: (error: any) => {
      arg.onShowError(error?.response?.data || "Fail to create user");
    },
  });

  const onSubmit: SubmitHandler<CreateUserFormData> = (data) => {
    mutate({
      ...data,
      permissions: data.isAdmin ? ["admin"] : [],
    });
  };

  return (
    <Dialog open={arg.isVisible} onClose={arg.close}>
      <DialogTitle>Create User</DialogTitle>
      <form onSubmit={handleSubmit(onSubmit)}>
        <DialogContent>
          <DialogContentText>
            To create a new user, please fill the form below. To create user
            with admin access, please check the "Set as admin" checkbox.
          </DialogContentText>
          <div className="h-4" />
          {[
            {
              fieldName: "name",
              label: "Name",
            },
            {
              fieldName: "username",
              label: "Username",
            },
            {
              fieldName: "password",
              label: "Password",
              type: "password",
            },
          ].map((fieldName) => (
            <>
              <TextField
                variant="outlined"
                fullWidth
                label={fieldName.label}
                type={fieldName.type || "text"}
                {...register(fieldName.fieldName as any)}
              />
              <div className="h-4" />
            </>
          ))}
          <FormControlLabel
            {...register("isAdmin")}
            control={<Checkbox />}
            label="Set as admin"
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={arg.close}>Cancel</Button>
          <Button type="submit">Register</Button>
        </DialogActions>
      </form>
    </Dialog>
  );
}
