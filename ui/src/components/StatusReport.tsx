import { CircularProgress, colors, List, ListItem, ListItemIcon, ListItemSecondaryAction, ListItemText, ListSubheader, makeStyles, Paper, Table, TableBody, TableCell, TableContainer, TableHead, TableRow } from '@material-ui/core';
import React, { useEffect, useState } from 'react'
import * as API from "../api";
import RouterIcon from '@material-ui/icons/Router';

import WbIncandescentIcon from '@material-ui/icons/WbIncandescent';
import CircularProgressWithLabel from './CircularProgressWithLabel';
import WatchLaterIcon from '@material-ui/icons/WatchLater';

import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

TimeAgo.addDefaultLocale(en)

/**---------------- */

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        maxWidth: 450,
        backgroundColor: theme.palette.background.paper,
        // color: theme.palette.text.primary,
    },
    text: {
        color: "blue",
    },
    table: {
        width: '100%',
        maxWidth: 500,
    },
    tableHead: {
        fontWeight: "bold"
    },
}));

/**---------------- */

export default function StatusReport() {
    const classes = useStyles();

    /*------------*/

    var loadTimeout: NodeJS.Timeout = null
    useEffect(() => {
        load();
        return () => {
            // Cleanup on unmount
            clearTimeout(loadTimeout);
            loadTimeout = null;
        }
    }, [])

    /*------------*/

    const [dataCollectionState, setDataCollectionState] = useState(null as API.DataCollectionStatus)
    const load = () => {
        API.getDataCollectionStatus().then(
            res => {
                // We need this as Go does not accept null for datetime
                if (new Date(res.LastExtractionTime).getFullYear() == 1) {
                    res.LastExtractionTime = null;
                }
                setDataCollectionState(res);
                loadTimeout = setTimeout(() => { load(); }, 2000);
            },
            err => { console.error(err); }
        )
    }

    /*------------*/

    const [lastNewSensorValues, setLastNewSensorValues] = useState(0)
    useEffect(() => {
        // Check if there is a change, load the data statistics
        if (lastNewSensorValues != dataCollectionState?.NewExtractedSensorValues) {
            setLastNewSensorValues(dataCollectionState?.NewExtractedSensorValues);
            loadTotalNumbers();
        }

    }, [dataCollectionState, lastNewSensorValues]);

    /*------------*/

    const [dataStatistics, setDataStatistics] = useState(null as API.DataCollectionStatistics)
    const loadTotalNumbers = () => {

        API.getDataCollectionStatistics().then(
            res => {
                setDataStatistics(res);
            },
            err => { console.error(err); }
        )
    }

    /*------------*/

    const niceNum = (n: number): string => {

        if (typeof n == "undefined") return "";
        if (n < 1000) return n.toString();
        if (n < 1000000) return Math.round(n / 1000).toString() + "K";
        if (n < 1000000000) return Math.round(n / 1000000).toString() + "M";
        return n.toString()
    }

    /*------------*/

    const timeAgo = new TimeAgo('en-US')

    /*------------*/

    return (
        <TableContainer component={Paper} className={classes.table}>
            <Table className={classes.table} aria-label="simple table">
                <TableHead >
                    <TableRow>
                        <TableCell className={classes.tableHead} align="center">Title</TableCell>
                        <TableCell className={classes.tableHead} align="right">New data</TableCell>
                        <TableCell className={classes.tableHead} align="right">Total</TableCell>
                    </TableRow>
                </TableHead>
                <TableBody>
                    <TableRow key={1}>
                        <TableCell component="th" scope="row">
                            <RouterIcon /> {"Devices"}
                        </TableCell>
                        <TableCell align="right" style={{ color: "blue" }}>
                            {dataCollectionState?.ChannelsRunning && <CircularProgress />}
                            {!dataCollectionState?.ChannelsRunning && <span title={dataCollectionState?.NewExtractedChannels?.toLocaleString()}>{niceNum(dataCollectionState?.NewExtractedChannels)}</span>}
                        </TableCell>
                        <TableCell align="right">
                            {dataStatistics?.totalChannels && <span title={dataStatistics?.totalChannels?.toLocaleString()}>{niceNum(dataStatistics?.totalChannels)}</span>}
                        </TableCell>
                    </TableRow>

                    <TableRow key={2}>
                        <TableCell component="th" scope="row">
                            <WbIncandescentIcon /> {"Sensor Values"}
                        </TableCell>
                        <TableCell align="right" style={{ color: "blue" }}>
                            {dataCollectionState?.SensorsRunning && <CircularProgressWithLabel size="55px" value={dataCollectionState?.SensorsProgress} />}
                            {!dataCollectionState?.SensorsRunning && <span color="primary" title={dataCollectionState?.NewExtractedSensorValues?.toLocaleString()}>{niceNum(dataCollectionState?.NewExtractedSensorValues)}</span>}
                        </TableCell>
                        <TableCell align="right">
                            {dataStatistics?.totalSensorValues && <span title={dataStatistics?.totalSensorValues?.toLocaleString()}>{niceNum(dataStatistics?.totalSensorValues)}</span>}
                        </TableCell>
                    </TableRow>

                    <TableRow key={3}>
                        <TableCell component="th" scope="row">
                            <WatchLaterIcon /> {"Last Extraction"}
                        </TableCell>
                        <TableCell align="right" colSpan={2}>
                            {dataCollectionState?.LastExtractionTime && timeAgo.format(new Date(dataCollectionState.LastExtractionTime))}
                        </TableCell>
                    </TableRow>
                </TableBody>
            </Table>
        </TableContainer>

    )
}
