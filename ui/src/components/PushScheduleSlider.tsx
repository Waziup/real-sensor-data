import React, { useState } from 'react';
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles';
import Typography from '@material-ui/core/Typography';
import Slider from '@material-ui/core/Slider';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            width: '95%',
        },
        margin: {
            height: theme.spacing(3),
        },
    }),
);

/*-----------*/

const marks = [
    { value: 0, realValue: 1, label: '1m', shortLabel: '1m', },
    { value: 5, realValue: 3, label: '3m', shortLabel: '3m', },
    { value: 10, realValue: 5, label: '5m', shortLabel: '5m', },
    { value: 15, realValue: 10, label: '10m', shortLabel: '10m', },
    { value: 25, realValue: 30, label: '30 mins', shortLabel: '30m', },
    { value: 35, realValue: 60, label: '1 hr', shortLabel: '1h', },
    { value: 45, realValue: 2 * 60, label: '2 hrs', shortLabel: '2h', },
    { value: 55, realValue: 3 * 60, label: '3 hrs', shortLabel: '3h', },
    { value: 65, realValue: 5 * 60, label: '5 hrs', shortLabel: '5h', },
    { value: 80, realValue: 24 * 60, label: '1 day', shortLabel: '1d', },
    { value: 90, realValue: 2 * 24 * 60, label: '2 days', shortLabel: '2d', },
    { value: 100, realValue: 3 * 24 * 60, label: '3 days', shortLabel: '3d', },
];

/*-----------*/

function valuetext(value: number) {
    return marks.find(o => o.value == value)?.label
}

function valueLabelFormat(value: number) {
    return marks.find(o => o.value == value)?.shortLabel
}

/*-----------*/

interface Props {
    onChange: (newValue: number) => void; // new value in minutes (realValue)
    defaultValue?: number; // realValue in minutes
    value?: number; // realValue in minutes
}

/*-----------*/

export default function PushScheduleSlider(props: Props) {
    const classes = useStyles();

    const handleChange = (event: React.ChangeEvent<{}>, newValue: number) => {
        const selected = marks.find(o => o.value == newValue)
        // console.log(newValue, selected?.realValue);
        props.onChange(selected?.realValue);
    }

    return (
        <div className={classes.root}>

            <Slider
                defaultValue={props.defaultValue ? marks.find(o => o.realValue == props.defaultValue)?.value : 10}
                getAriaValueText={valuetext}
                valueLabelFormat={valueLabelFormat}
                aria-labelledby="schedule-slider"
                step={null}
                max={marks[marks.length - 1].value}
                scale={(x) => x}
                valueLabelDisplay="auto"
                value={props.value && marks.find(o => o.realValue == props.value)?.value}
                marks={marks}
                onChange={handleChange}
            />
        </div>
    );
}
