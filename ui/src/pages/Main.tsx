import React, { useEffect, useState } from 'react';
import { AppBar, Button, Dialog, DialogActions, DialogContent, DialogContentText, DialogTitle, IconButton, makeStyles, Toolbar, Typography } from '@material-ui/core';
import EqualizerIcon from '@material-ui/icons/Equalizer';

// import ontologies from "../../data/ontologies/ontologies.json";
import StatusReport from '../layouts/StatusReport';
import SensorSearch from '../layouts/SensorSearch';
import Sensor from '../layouts/Sensor';
import BackToTop from '../components/BackToTop';

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

  useEffect(() => {
    const tmpData = { id: 284943, title: "Air Temp (F)", subtitle: "Daniels Island Tide and Meteorological Station" };
    handleSearchResultClick(tmpData);
    // return () => {}
  }, [])

  /**------------- */

  const [statsDlgOpen, setStatisticsDialogue] = useState(false)
  const handleStatisticsOpen = () => { setStatisticsDialogue(true) }
  const handleStatisticsClose = () => { setStatisticsDialogue(false) }

  /**------------- */



  /**------------- */

  const [sensorDlgOpen, setSensorDlgOpen] = useState(false)
  const [sensorProps, setSensorProps] = useState(null)
  const handleSearchResultClick = (dataRow: any) => {
    setSensorDlgOpen(true)
    setSensorProps({
      channel_id: dataRow.id,
      name: dataRow.title,
    })
  }
  const handleSensorDlgClose = () => { setSensorDlgOpen(false) }

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

      {/* ------------------------- */}

      <Toolbar id="back-to-top-anchor" />

      <SensorSearch onSearchResultClick={handleSearchResultClick} />

      <BackToTop topAnchorId="back-to-top-anchor" />

      {/* ------------------------- */}

      <Dialog
        open={sensorDlgOpen}
        onClose={handleSensorDlgClose}
        aria-labelledby="sensor-dialog-title"
        aria-describedby="sensor-dialog-description"
        fullScreen={true}
      >
        <DialogTitle id="sensor-dialog-title">{"Sensor Details"}</DialogTitle>
        <DialogContent>
          {sensorDlgOpen && <Sensor {...sensorProps} />}
        </DialogContent>
        <DialogActions>
          <Button onClick={handleSensorDlgClose} color="primary" autoFocus>
            Close
          </Button>
        </DialogActions>
      </Dialog>

      {/* ------------------------- */}

      <Dialog
        open={statsDlgOpen}
        onClose={handleStatisticsClose}
        aria-labelledby="statistics-dialog-title"
        aria-describedby="statistics-dialog-description"
      >
        <DialogTitle id="statistics-dialog-title">{"Data Collection Statistics"}</DialogTitle>
        <DialogContent>
          <StatusReport />
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