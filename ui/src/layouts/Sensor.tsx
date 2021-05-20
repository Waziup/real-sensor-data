import { CircularProgress, Grid, makeStyles, Paper, Table, TableBody, TableCell, TableContainer, TableRow } from '@material-ui/core';
import React, { useEffect, useState } from 'react'
import * as API from "../api";
import Chart, { TimeSeriesDataType } from '../components/Chart';

/**---------------- */


import TimeAgo from 'javascript-time-ago'
import en from 'javascript-time-ago/locale/en'

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
    channel_id: number;
    name: string;
}

/**---------------- */

export default function Sensor(props: Props) {

    const classes = useStyles();
    const timeAgo = new TimeAgo('en-US')

    /**------------------ */

    useEffect(() => {
        loadChannel();
        loadSensorValues();
        // return () => {}
    }, [])

    /**------------------ */

    const [channel, setChannel] = useState(null as API.ChannelRow)
    const loadChannel = () => {

        API.getChannel(props.channel_id).then(
            res => { setChannel(res); },
            err => { console.error(err); }
        )
    }

    /**------------------ */

    const [totalEntries, setTotalEntries] = useState(0)
    const [chartValues, setSensorValues] = useState(null as TimeSeriesDataType[])
    const loadSensorValues = () => {
        API.getSensorValues(props.name, props.channel_id).then(
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
            { title: "Sensor name", value: props.name },
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


        // `https://www.google.com/maps/search/52.759723,-1.236892`
    }

    /**------------------ */

    if (chartValues === null) {
        return (<Grid container justify="center">
            <CircularProgress disableShrink />
        </Grid>)
    }

    /**------------------ */

    return (
        <div>
            <TableContainer component={Paper} className={classes.table}>
                <Table className={classes.table} aria-label="simple table">
                    <TableBody>
                        {tableInfo && tableInfo.map((row: any, index: number) =>
                            <TableRow key={index}>
                                {row.title &&
                                    <TableCell component="td" width={250} scope="row">{row.title}</TableCell>
                                }
                                <TableCell align="left" colSpan={row.title ? 1 : 2} style={{ color: "#06A" }}>{row.value}</TableCell>
                            </TableRow>
                        )}

                    </TableBody>
                </Table>
            </TableContainer>
            {chartValues && <Chart data={chartValues} />}
        </div>
    )
}
