import React, { useState } from "react";
import {
  Button,
  Paper,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
} from "@mui/material";
import { useMutation, useQuery } from "react-query";
import { deleteUser, users } from "../../api/auth";
import DeleteIcon from "@mui/icons-material/Delete";
import AlertDialog from "../../components/dialog";
import CreateUser from "./create-user";
import { useSession } from "../../context/session-context";

const UsersPage = () => {
  const { session } = useSession();

  const { data, refetch } = useQuery("users", () => users());
  const { mutate: mutateDelete } = useMutation("deleteUser", deleteUser, {
    onSuccess: () => refetch(),
    onSettled: () => setShowDeleteUserDialog([false, ""]),
  });

  const [showDeleteUserDialog, setShowDeleteUserDialog] = React.useState<
    [boolean, string]
  >([false, ""]);

  const [showCreateUserDialog, setShowCreateUserDialog] = React.useState<
    [boolean, string]
  >([false, ""]);

  const [showRegister, setShowRegister] = useState(false);

  return (
    <>
      <div>
        {session?.permissions?.includes("admin") && (
          <Button
            variant="contained"
            sx={{
              backgroundColor: "#3f3100",
              color: "#fff0bf",
              "&:hover": {
                backgroundColor: "#fff0bf",
                color: "#3f3100",
              },
            }}
            onClick={() => setShowRegister(true)}
          >
            New
          </Button>
        )}
        <TableContainer component={Paper} className="mt-4">
          <Table>
            <TableHead>
              <TableRow>
                <TableCell>Name</TableCell>
                <TableCell>Permissions</TableCell>
                <TableCell />
              </TableRow>
            </TableHead>
            <TableBody>
              {data?.map((row) => (
                <TableRow key={row.userId}>
                  <TableCell>{row.name}</TableCell>
                  <TableCell>{row.permissions?.join(",")}</TableCell>
                  {session?.userId !== row.userId ? (
                    <TableCell
                      onClick={() =>
                        setShowDeleteUserDialog([true, row.userId])
                      }
                      style={{
                        width: "10px",
                        whiteSpace: "normal",
                        wordWrap: "break-word",
                      }}
                    >
                      <DeleteIcon />
                    </TableCell>
                  ) : (
                    <TableCell />
                  )}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      </div>
      <AlertDialog
        isOpen={showDeleteUserDialog[0]}
        handleClose={() => setShowDeleteUserDialog([false, ""])}
        title={"Delete user"}
        description={"Are you sure you want to delete this user?"}
        buttons={[
          <Button
            key="cancel"
            onClick={() => setShowDeleteUserDialog([false, ""])}
            color="primary"
          >
            Cancel
          </Button>,
          <Button
            key="logout"
            onClick={() => {
              mutateDelete(showDeleteUserDialog[1]);
            }}
            color="primary"
            autoFocus
          >
            Delete
          </Button>,
        ]}
      />
      <AlertDialog
        isOpen={showCreateUserDialog[0]}
        handleClose={() => setShowCreateUserDialog([false, ""])}
        title={"Create user failed"}
        description={showCreateUserDialog[1]}
        buttons={[
          <Button
            key="cancel"
            onClick={() => setShowCreateUserDialog([false, ""])}
            color="primary"
          >
            OK
          </Button>,
        ]}
      />
      <CreateUser
        isVisible={showRegister}
        close={() => setShowRegister(false)}
        onShowError={(message) => setShowCreateUserDialog([true, message])}
        onUserCreated={() => refetch()}
      />
    </>
  );
};

export default UsersPage;
