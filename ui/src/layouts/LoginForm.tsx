import React, { useState, Fragment, useEffect } from "react";
// import { TimeComp } from "../Time";

import LockOpenIcon from '@material-ui/icons/LockOpen';
import Container from '@material-ui/core/Container';
import { Alert } from '@material-ui/lab';
import * as API from "../api";

import {
    Divider,
    Card,
    CardContent,
    makeStyles,
    Grid,
    List,
    FormGroup,
    TextField,
    Button,
    CardActions,
    Grow,
    LinearProgress,
    ListItem,
    ListItemText
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
    onSuccess: (user: API.UserType) => void;
    onFailure?: (err: string) => void;
}

/**------------ */

export default function LoginForm(props: Props) {
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
    const [loginErr, setLoginErr] = useState(false);
    const loginCheck = (event: any) => {
        event.preventDefault();
        setChecking(true);
        API.login(event.target.username.value, event.target.password.value).then(
            (res) => {
                setChecking(false);
                setLoginErr(false);
                console.log("Token", res);
                setMsg("Login Success! [ Redirecting... ]");
                //Redirecting...

                props.onSuccess({
                    username: event.target.username.value,
                    tokenHash: res
                });

            },
            (error) => {
                setChecking(false);
                setLoginErr(true);
                console.log(error);
                setMsg("Invalid credentials!");
                if (props.onFailure) {
                    props.onFailure(error)
                }
            }
        );
    };


    return (

        <div className={classes.paper}>
            <Card className={classes.root}>
                <Grow in={checking}>
                    <LinearProgress />
                </Grow>
                <List dense={true}>
                    <ListItem>
                        {/* <img className={classes.logo} src={wazigateLogo} /> */}
                        <ListItemText
                            primary="Login to with your Waziup Credentials"
                        />
                    </ListItem>
                </List>
                <Divider />
                <CardContent>
                    <form noValidate onSubmit={loginCheck}>
                        <FormGroup>

                            <TextField
                                // variant="outlined"
                                margin="normal"
                                required
                                fullWidth
                                id="username"
                                label="Username"
                                name="username"
                                autoComplete="username"
                                autoFocus
                            />
                            <TextField
                                // variant="outlined"
                                margin="normal"
                                required
                                fullWidth
                                name="password"
                                label="Password"
                                type="password"
                                id="password"
                                autoComplete="current-password"
                            />
                            <CardActions>
                                <Button
                                    type="submit"
                                    fullWidth
                                    disabled={checking}
                                    color="primary"
                                    variant="contained"
                                    startIcon={<LockOpenIcon />}
                                >
                                    Login
                </Button>
                            </CardActions>

                        </FormGroup>
                    </form>

                    <Divider />
                    <br />

                    {msg != "" && (<Alert severity={loginErr ? "error" : "success"}>{msg}</Alert>)}

                    <Grid container>
                        <Grid item xs>
                            If you do not have an account please head over
                            to <a target="_blank" href="https://dashboard.waziup.io/">https://dashboard.waziup.io/</a>
                            create an account.
                        </Grid>
                    </Grid>
                </CardContent>
            </Card>
        </div>
    );
}