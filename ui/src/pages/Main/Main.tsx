import React, { useEffect, useState } from 'react';
import { AppBar, Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, makeStyles, Toolbar, Typography } from '@material-ui/core';
import EqualizerIcon from '@material-ui/icons/Equalizer';

// import ontologies from "../../data/ontologies/ontologies.json";
import StatusReport from '../../components/StatusReport';
import SensorSearch from '../../components/SensorSearch';

/*---------------------*/

const useStyles = makeStyles((theme) => ({
  root: {
    '& > *': {
      margin: theme.spacing(1),
      width: '90%',
      display: 'flex'
    },
  },
  button: {
    width: '6rem',
    margin: theme.spacing(1),
  },
  box: {
    border: 'solid 1px #CCC',
    borderRadius: theme.spacing(1),
    margin: theme.spacing(1),
    padding: theme.spacing(1),
    textAlign: "left",
  },
  menuButton: {
    marginRight: theme.spacing(2),
  },
  title: {
    flexGrow: 1,
  },
  fieldset: {
    justifyContent: "center",
    padding: "1rem"
  },
  numInput: {
    width: '5rem',
    marginLeft: '.5rem'
  },
}));

/*---------------------*/

export default function Main() {
  const classes = useStyles();

  /**------------- */

  const [statsDlgOpen, setStatisticsDialogue] = useState(false)
  const handleStatisticsOpen = () => { setStatisticsDialogue(true) }
  const handleStatisticsClose = () => { setStatisticsDialogue(false) }

  /**------------- */


  return (
    <div className="Main">
      <AppBar position="static">
        <Toolbar>
          {/* <IconButton edge="start" className={classes.menuButton} color="inherit" aria-label="menu">
            <MenuIcon />
          </IconButton> */}
          <Typography variant="h6" className={classes.title}>
            Sensor Data Simulator
           </Typography>
          {/* <Button color="inherit">Login</Button> */}
          <IconButton
            edge="end"
            aria-label="Statistics"
            aria-haspopup="true"
            onClick={handleStatisticsOpen}
            color="inherit"
          >
            <EqualizerIcon />
          </IconButton>
        </Toolbar>
      </AppBar>

      <SensorSearch />


      <Dialog
        open={statsDlgOpen}
        onClose={handleStatisticsClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
      >
        <DialogTitle id="alert-dialog-title">{"Data Collection Statistics"}</DialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            <StatusReport />
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleStatisticsClose} color="primary" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>

    </div >
  );
}

/**---------------- */