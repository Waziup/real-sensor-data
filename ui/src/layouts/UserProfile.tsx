import React, { useState, Fragment, useEffect } from "react";
// import { TimeComp } from "../Time";

import ExitToAppIcon from '@material-ui/icons/ExitToApp';

import { Alert } from '@material-ui/lab';
import * as API from "../api";

import {
    Card,
    CardContent,
    makeStyles,
    Grid,
    Button,
    CardActions,
    Grow,
    LinearProgress,
} from '@material-ui/core';


const useStyles = makeStyles((theme) => ({
    root: {
        maxWidth: 600,
        // maxWidth: "calc(100% - 32px)",
        display: "inline-block",
        verticalAlign: "top",
    },
    name: {
        cursor: "text",
        '&:hover': {
            "text-decoration": "underline",
        },
    },
    icon: {
        width: "40px",
        height: "40px",
    },
    logo: {
        display: "inline-flex",
        height: "2rem",
        marginRight: 16,
    },
    expand: {
        transform: 'rotate(0deg)',
        marginLeft: 'auto',
        transition: theme.transitions.create('transform', {
            duration: theme.transitions.duration.shortest,
        }),
    },
    expandOpen: {
        transform: 'rotate(180deg)',
    },
    value: {
        float: "right",
        flexGrow: 0,
        marginLeft: "1.5em",
    },
    wrapper: {
        position: "relative",
    },
    progress: {
        color: "#4caf50",
        display: "inline",
    },
    paper: {
        marginTop: theme.spacing(8),
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
    }
}));

/**------------ */

interface Props {
    onLogout?: () => void;
    user: API.UserType;
}

/**------------ */

export default function UserProfile(props: Props) {
    const classes = useStyles();

    /*------------ */
    // Run stuff on load
    useEffect(
        () => {
            // Nothing to run yet
        },
        [] /* This makes it to run only once*/
    );

    /*------------ */
    const [msg, setMsg] = useState("");
    const [checking, setChecking] = useState(false);
    const [loginErr, setErr] = useState(false);
    const logout = (event: any) => {
        // event.preventDefault();
        if (!confirm("Do you really want to logout?")) return;
        setChecking(true);

        API.logout().then(
            (res) => {
                setChecking(false);
                setErr(false);
                setMsg("Logout Success!");
                //Redirecting...

                if (props.onLogout) props.onLogout();

            },
            (error) => {
                setChecking(false);
                setErr(true);
                console.log(error);
                setMsg("Failed to logout!");
            }
        );
    };


    return (

        <div className={classes.paper}>
            <Card className={classes.root}>
                <Grow in={checking}>
                    <LinearProgress />
                </Grow>
                {/* <List dense={true}>
                    <ListItem>
                        <ListItemText
                            primary="User profile"
                        />
                    </ListItem>
                </List> */}
                <CardContent>
                    <Grid container style={{ margin: 20 }}>
                        <Grid item>
                            Logged in as: <span style={{ color: "blue" }}>{props.user.username}</span>
                        </Grid>
                    </Grid>
                    <CardActions>
                        <Button
                            fullWidth
                            disabled={checking}
                            color="secondary"
                            variant="contained"
                            startIcon={<ExitToAppIcon />}
                            onClick={logout}
                        >Logout</Button>
                    </CardActions>

                    {msg != "" && (<Alert severity={loginErr ? "error" : "success"}>{msg}</Alert>)}

                </CardContent>
            </Card>
        </div>
    );
}