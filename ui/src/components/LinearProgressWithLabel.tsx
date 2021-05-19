import React from 'react';
import PropTypes from 'prop-types';
import { makeStyles } from '@material-ui/core/styles';
import LinearProgress from '@material-ui/core/LinearProgress';
import Typography from '@material-ui/core/Typography';
import Box from '@material-ui/core/Box';



type Props = {
    /**
     * The value of the progress indicator for the determinate and buffer variants.
     * Value between 0 and 100.
     */
    value: number,
};

const useStyles = makeStyles({
    root: {
        width: '100%',
    },
});

export default function LinearWithValueLabel({ value }: Props) {
    const classes = useStyles();
    return (
        <div className={classes.root}>
            <Box display="flex" alignItems="center">
                <Box width="100%" mr={1}>
                    <LinearProgress variant="determinate" value={value} />
                </Box>
                <Box minWidth={35}>
                    <Typography variant="body2" color="textSecondary">{`${Math.round(
                        value,
                    )}%`}</Typography>
                </Box>
            </Box>
        </div>
    );
}
