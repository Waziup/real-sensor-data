import React from "react";
import { Zoom, useScrollTrigger, Fab, Theme, makeStyles, createStyles, Color } from "@material-ui/core"
import { KeyboardArrowUp } from "@material-ui/icons";


const useStyles = makeStyles((theme: Theme) =>
    createStyles({
        root: {
            position: 'fixed',
            bottom: theme.spacing(2),
            right: theme.spacing(2),
        },
    }),
);

/*---------*/

interface Props {
    topAnchorId: string;
    color?: Color;
    window?: () => Window;
}

export default function BackToTop(props: Props) {

    const classes = useStyles()

    const trigger = useScrollTrigger({
        target: props.window ? props.window() : undefined,
        disableHysteresis: true,
        threshold: 100,
    });

    /**-------- */

    const handleClick = (event: React.MouseEvent<HTMLDivElement>) => {
        const anchor = ((event.target as HTMLDivElement).ownerDocument || document).querySelector(`#${props.topAnchorId}`)
        if (anchor) {
            anchor.scrollIntoView({ behavior: "smooth", block: "center" })
        }
    }

    /**-------- */
    return (
        <Zoom in={trigger}>
            <div onClick={handleClick} role="presentation" className={classes.root} >
                <Fab color={props.color || "secondary" as any} size="large" aria-label="scroll back to top">
                    <KeyboardArrowUp />
                </Fab>
            </div>
        </Zoom>
    )
}

/*---------*/
