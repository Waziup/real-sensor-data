import React from 'react';
import { makeStyles, Theme, createStyles } from '@material-ui/core/styles';
import Paper from '@material-ui/core/Paper';
import InputBase from '@material-ui/core/InputBase';
import IconButton from '@material-ui/core/IconButton';
import SearchIcon from '@material-ui/icons/Search';

const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            padding: '2px 4px',
            display: 'flex',
            alignItems: 'center',
            width: '100%',
        },
        input: {
            marginLeft: theme.spacing(1),
            flex: 1,
        },
        iconButton: {
            padding: 10,
        },
        divider: {
            height: 28,
            margin: 4,
        },
    }),
);

/**--------------- */

type Props = {
    placeholder?: string;
    label?: string;
    onChange?: React.ChangeEventHandler<HTMLTextAreaElement | HTMLInputElement>;
    onSubmit?: React.FormEventHandler<HTMLFormElement | HTMLDivElement>;
    width?: number;
}


/**--------------- */

export default function SearchBar(props: Props) {
    const classes = useStyles();

    return (
        <Paper component="form" className={classes.root} onSubmit={props.onSubmit}>
            <InputBase
                className={classes.input}
                placeholder={props.placeholder || "Search"}
                inputProps={{ 'aria-label': "Search" }}
                onChange={props.onChange}
            />
            <IconButton type="submit" className={classes.iconButton} aria-label={props.label || "Search"}>
                <SearchIcon />
            </IconButton>
        </Paper>
    );
}