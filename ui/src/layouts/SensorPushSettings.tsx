import { Button, CircularProgress, createStyles, Grid, makeStyles, Switch, Theme, Typography } from '@material-ui/core';
import React, { useEffect, useState } from 'react'
import TextField from '@material-ui/core/TextField';
import Autocomplete from '@material-ui/lab/Autocomplete';
import SaveIcon from '@material-ui/icons/Save';
import AddBoxIcon from '@material-ui/icons/AddBox';
import ClearIcon from '@material-ui/icons/Clear';
import DeleteForeverIcon from '@material-ui/icons/DeleteForever';

import * as API from "../api";
import LoginForm from './LoginForm';
import PushScheduleSlider from '../components/PushScheduleSlider';
import DataTablePushSettings from './DataTablePushSettings';


/**---------- */

interface Props {
    sensorId: number;
}

/**---------- */

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            flexGrow: 1,
        },
        paper: {
            padding: theme.spacing(2),
            textAlign: 'center',
            color: theme.palette.text.secondary,
        },
    }),
);

/**---------- */

export default function SensorPushSettings(props: Props) {

    const classes = useStyles();

    /**--------------- */

    useEffect(() => {

        loadDevices();
        // return () => {}
    }, [])


    /**--------------- */

    const [originalTimestamp, setOriginalTimestamp] = useState(false)
    const handleOriginalTimestamp = (event: React.ChangeEvent<HTMLInputElement>) => {
        setOriginalTimestamp(event.target.checked);
    }
    /**--------------- */

    const [activePush, setActivePush] = useState(true)
    const handleActiveChange = (event: React.ChangeEvent<HTMLInputElement>) => {
        setActivePush(event.target.checked);
    }

    /**--------------- */

    const [editMode, setEditMode] = useState(false)
    const [savingPushSettings, setSavingPushSettings] = useState<boolean>(false)
    const savePushSettings = () => {

        if (!selectedSensor || !selectedSensor?.id) { return }

        const data: API.SensorPushSettings = {
            target_device_id: selectedSensor.devId,
            target_sensor_id: selectedSensor.id,
            active: activePush,
            id: editMode ? recordId : 0, //default:0 means insert a new record
            push_interval: pushInterval,
            use_original_time: originalTimestamp,
        }

        setSavingPushSettings(true);
        API.savePushSettings(data, props.sensorId).then(
            res => {
                let tmp = userDevices
                setUserDevices(null);
                setUserDevices(tmp); // Just to make the table refresh ;)

                setErr(null);
            },
            err => {
                setErr(err);
                console.error(err);
            }
        ).finally(() => setSavingPushSettings(false))
    }

    /**--------------- */

    const [userDevices, setUserDevices] = useState<API.SensorType[]>(null)
    const [err, setErr] = useState<API.HttpError>(null)
    const [loading, setLoading] = useState<boolean>(false)
    const loadDevices = () => {
        setLoading(true);
        API.getUserDevices().then(
            res => {

                // Refine the data
                let devList: API.SensorType[] = new Array()
                if (res) {
                    for (let dev of res) {
                        for (let sensor of dev.sensors) {
                            devList.push({
                                devName: dev.name,
                                devId: dev.id,
                                name: sensor.name,
                                id: sensor.id,
                                title: dev.name + ": " + sensor.name
                            });
                        }
                    }
                }

                setUserDevices(devList);
                setErr(null);
            },
            err => {
                setErr(err);
                console.error(err);
            }
        ).finally(() => setLoading(false))
    }


    /**--------------- */

    const [selectedSensor, setSelectedSensor] = useState<API.SensorType>(null)
    const handleSensorSelect = (e: any, val: API.SensorType) => {
        setSelectedSensor(val);
    }

    const [pushInterval, setPushInterval] = useState(10)
    const handleIntervalSelect = (value: number) => {
        setPushInterval(value);
    }

    /**--------------- */

    let options: any;
    if (userDevices !== null) {
        options = userDevices.map((option: any) => {
            const firstLetter = option.title[0].toUpperCase();
            return {
                firstLetter: /[0-9]/.test(firstLetter) ? '0-9' : firstLetter,
                ...option,
            };
        });
    }

    /**--------------- */

    const [recordId, setRecordId] = useState(0)
    const handleTableRowClick = (data: API.SensorPushSettings) => {

        setEditMode(true);
        setPushInterval(data.push_interval);
        setActivePush(data.active);
        setRecordId(data.id);
        setOriginalTimestamp(data.use_original_time == true);

        /**--------- */

        const devName = userDevices.find(o => o.devId == data.target_device_id)?.devName || data.target_device_id;
        const sensorName = userDevices.find(o => o.id == data.target_sensor_id)?.name || data.target_sensor_id;
        // const firstLetter = devName[0].toUpperCase();

        setSelectedSensor({
            // firstLetter: /[0-9]/.test(firstLetter) ? '0-9' : firstLetter,
            devName: devName,
            devId: data.target_device_id,
            name: sensorName,
            id: data.target_sensor_id,
            title: devName + ": " + sensorName
        });

        /**--------- */
    }

    /**--------------- */

    const [deleting, setDeleting] = useState(false)
    const deleteRecord = () => {
        if (!confirm("Do you really want to delete this setting?")) return;
        setDeleting(true)
        API.deletePushSettings(props.sensorId, recordId).then(
            res => {
                let tmp = userDevices
                setUserDevices(null);
                setUserDevices(tmp); // Just to make the table refresh ;)
                setErr(null);
            },
            err => {
                setErr(err);
                console.error(err);
            }
        ).finally(() => setDeleting(false))
    }

    /**--------------- */


    const resetForm = () => {

        setEditMode(false);
        setSelectedSensor(null);
        setPushInterval(10);
        setActivePush(true);
        setRecordId(0);
        setOriginalTimestamp(false);
    }

    /**--------------- */

    // If authorization failed
    if (err && err.status == 401) {
        return <LoginForm onSuccess={() => { loadDevices(); }} />
    }

    /**--------------- */


    if (loading || userDevices === null) {
        return (<Grid container justify="center">
            <CircularProgress disableShrink />
        </Grid>)
    }

    /**--------------- */
    return (
        <Grid container spacing={3}>

            <Grid item xs={3}>
                <Typography gutterBottom>Target Sensor</Typography>
            </Grid>
            <Grid item xs={9}>
                <Autocomplete
                    id="user-devices"
                    options={options.sort((a: any, b: any) => -b.firstLetter.localeCompare(a.firstLetter))}
                    groupBy={(option: any) => option.devName}
                    getOptionLabel={(option: any) => option.title}
                    style={{ width: 600 }}
                    renderInput={(params) => <TextField {...params} label="Devices and Sensors" variant="outlined" />}
                    onChange={handleSensorSelect}
                    value={selectedSensor}
                />
            </Grid>

            <Grid item xs={3}>
                <Typography gutterBottom>Push Interval</Typography>
            </Grid>
            <Grid item xs={9}>
                <PushScheduleSlider value={pushInterval} defaultValue={pushInterval} onChange={handleIntervalSelect} />
            </Grid>

            <Grid item xs={3}>
                <Typography gutterBottom>Activate</Typography>
            </Grid>
            <Grid item xs={9}>
                <Switch
                    checked={activePush}
                    onChange={handleActiveChange}
                    color="primary"
                    name="activePush"
                    inputProps={{ 'aria-label': 'primary checkbox' }}
                />
            </Grid>

            <Grid item xs={3}>
                <Typography gutterBottom>Use Original Timestamp</Typography>
            </Grid>
            <Grid item xs={9}>
                <Switch
                    checked={originalTimestamp}
                    onChange={handleOriginalTimestamp}
                    color="primary"
                    name="originalTimestamp"
                    inputProps={{ 'aria-label': 'primary checkbox' }}
                />
            </Grid>

            <Grid item xs={12}><br /></Grid>

            <Grid item xs={3}></Grid>
            <Grid item xs={3}>
                {editMode && <Button
                    fullWidth
                    color="inherit"
                    variant="contained"
                    startIcon={<ClearIcon />}
                    onClick={() => resetForm()}
                >Cancel
                </Button>}
            </Grid>
            <Grid item xs={3}>
                {editMode && <Button
                    fullWidth
                    disabled={savingPushSettings}
                    color="secondary"
                    variant="contained"
                    startIcon={<DeleteForeverIcon />}
                    onClick={() => deleteRecord()}
                >Delete
                {deleting && <CircularProgress color="primary" />}
                </Button>}
            </Grid>
            <Grid item xs={3}>
                <Button
                    type="submit"
                    fullWidth
                    disabled={savingPushSettings}
                    color="primary"
                    variant="contained"
                    startIcon={editMode ? <SaveIcon /> : <AddBoxIcon />}
                    onClick={() => savePushSettings()}
                >
                    {editMode ? "Save" : "Add"}
                    {savingPushSettings && <CircularProgress color="secondary" />}
                </Button>
            </Grid>

            <Grid item xs={12}>
                <br /><br />
                <DataTablePushSettings sensorId={props.sensorId} userDevices={userDevices} onRowClick={handleTableRowClick} />
            </Grid>
        </Grid>
    );

}

