import { AppBar, Box, CircularProgress, Divider, Grid, makeStyles, Paper, Tab, Table, TableBody, TableCell, TableContainer, TableRow, Tabs } from '@material-ui/core';
import React, { useEffect, useState } from 'react'
import InfoIcon from '@material-ui/icons/Info';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import CloudUploadIcon from '@material-ui/icons/CloudUpload';
import Alert from '@material-ui/lab/Alert';

import * as API from "../api";
import Chart, { TimeSeriesDataType } from '../components/Chart';
import SensorPushSettings from './SensorPushSettings'

/**---------------- */

import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'
let timeAgo: TimeAgo
try {
    TimeAgo.addDefaultLocale(en)
} catch (e) { //console.warn(e) 
}


/**---------------- */

const useStyles = makeStyles((theme) => ({
    root: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
        // color: theme.palette.text.primary,
        flexGrow: 1,
    },
    tabs: {
        width: '100%',
        backgroundColor: theme.palette.background.paper,
    },
    tabPanels: {
        border: '1px solid #EEE',
        borderRadius: 5
    },
    text: {
        color: "blue",
    },
    table: {
        width: '100%',
    },
    tableHead: {
        fontWeight: "bold"
    },
}));

/**---------------- */

interface Props {
    id: number;
}

/**---------------- */

export default function Sensor(props: Props) {

    const classes = useStyles();
    // const theme = useTheme();


    /**------------------ */

    useEffect(() => {
        timeAgo = new TimeAgo('en-US')

        loadChannel();
        loadSensorValues();
        // return () => {}
    }, [])

    /**------------------ */

    const [err, setErr] = useState<string>(null)
    const [channel, setChannel] = useState<API.ChannelRow>(null)
    const [sensorName, setSensorName] = useState<string>(null)
    const loadChannel = () => {

        API.getSensor(props.id).then(
            res => {
                setSensorName(res.name);
                API.getChannel(res.channel_id).then(
                    res => { setChannel(res); },
                    err => { console.error(err); }
                )
            },
            err => {
                setErr(err);
                console.error(err);
            }
        )
    }

    /**------------------ */



    /**------------------ */

    const [totalEntries, setTotalEntries] = useState(0)
    const [chartValues, setSensorValues] = useState<TimeSeriesDataType[]>(null)
    const loadSensorValues = () => {
        API.getSensorValues(props.id, 1).then(
            res => {
                setTotalEntries(res.pagination.total_entries);

                let chartData: TimeSeriesDataType[] = new Array()
                if (res.rows) {
                    for (let row of res.rows) {
                        chartData.push({
                            time: row.created_at,
                            value: row.value as any
                        })
                    }
                }
                setSensorValues(chartData);
            },
            err => { console.error(err); }
        )
    }

    /**------------------ */

    // Preparing the info to show
    let tableInfo: any = null;
    if (channel !== null && chartValues !== null) {

        tableInfo = [
            { title: "Sensor name", value: sensorName },
            { title: "Device name", value: channel.name },
            { title: "", value: channel.description },
            { title: "Created Time", value: new Date(channel.created_at).toLocaleString() },
            { title: "Last Activity", value: chartValues.length ? timeAgo.format(new Date(chartValues[0]?.time)) : null },
        ];

        if (channel.latitude && channel.longitude) {
            tableInfo.push({ title: "Location", value: <a target="_blank" href={`https://www.google.com/maps/search/${channel.latitude},${channel.longitude}`}>{channel.latitude + " , "}{channel.longitude}</a> })
        }

        if (channel.url) {
            tableInfo.push({ title: "URL", value: <a target="_blank" href={channel.url}>{channel.url}</a> })
        }

        tableInfo.push({ title: "Total value entries", value: totalEntries.toLocaleString() })

    }

    /**------------------ */

    const [tabValue, setTabValue] = useState(0)
    const handleTabChange = (event: React.ChangeEvent<{}>, newValue: number) => {
        setTabValue(newValue);
    };

    /**------------------ */

    if (err) {
        return <Alert severity="error">{err}</Alert>
    }

    /**------------------ */


    if (tableInfo === null) {
        return (<Grid container justify="center">
            <CircularProgress disableShrink />
        </Grid>)
    }

    /**------------------ */

    return (
        <div className={classes.root}>
            <AppBar position="static" color="default">
                <Tabs
                    value={tabValue}
                    onChange={handleTabChange}
                    variant="fullWidth" // | scrollable
                    // scrollButtons="on"
                    indicatorColor="primary"
                    textColor="primary"
                    aria-label="sensor details tabs"
                >
                    <Tab label="Details" icon={<InfoIcon />} {...a11yProps(0)} />
                    <Tab label="Chart" icon={<ShowChartIcon />} {...a11yProps(1)} />
                    <Tab label="Push Settings" icon={<CloudUploadIcon />} {...a11yProps(2)} />
                </Tabs>
            </AppBar>
            <Box className={classes.tabPanels}>
                <TabPanel value={tabValue} index={0} >
                    <TableContainer >
                        <Table className={classes.table} aria-label="simple table">
                            <TableBody>
                                {tableInfo && tableInfo.map((row: any, index: number) =>
                                    <TableRow key={index}>{row.title != "" &&
                                        <TableCell component="td" width={250} scope="row">{row.title}</TableCell>
                                    }
                                        <TableCell align="left" colSpan={row.title ? 1 : 2} style={{ color: "#06A" }}>
                                            {row.value}
                                        </TableCell>
                                    </TableRow>
                                )}
                            </TableBody>
                        </Table>
                    </TableContainer>
                </TabPanel>
                <TabPanel value={tabValue} index={1}>
                    {chartValues && <Chart data={chartValues} />}
                </TabPanel>
                <TabPanel value={tabValue} index={2}>
                    <Grid container spacing={3}>
                        <Grid item xs={3}>Source Sensor</Grid>
                        <Grid item xs={9}>{channel?.name && sensorName ? (channel?.name + " - " + sensorName) : "..."}</Grid>
                    </Grid>
                    <br />
                    <Divider />
                    <br />
                    <SensorPushSettings sensorId={props.id} />
                </TabPanel>
            </Box>
        </div>
    )
}

/**----------------- */

interface TabPanelProps {
    children?: React.ReactNode;
    index: any;
    value: any;
}

function TabPanel(props: TabPanelProps) {
    const { children, value, index, ...other } = props;

    return (
        <div
            role="tabpanel"
            hidden={value !== index}
            id={`tabpanel-${index}`}
            aria-labelledby={`tab-${index}`}
            {...other}
        >
            {value === index && (
                <Box p={3}>
                    {children}
                </Box>
            )}
        </div>
    );
}

/**----------------- */

function a11yProps(index: any) {
    return {
        id: `tab-${index}`,
        'aria-controls': `tabpanel-${index}`,
    };
}

/**----------------- */