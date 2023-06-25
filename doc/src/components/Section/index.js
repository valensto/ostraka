import React from "react"
import style from "./styles.module.css"
import clsx from "clsx"

import {Title} from "./Title"
import {Subtitle} from "./Subtitle"

export const Section = ({
                            fullWidth,
                            children,
                            odd,
                            accent,
                            row,
                            noGap,
                            center,
                            className = "",
                        }) => (
    <div
        className={clsx(
            style.root,
            {
                [style.odd]: odd,
                [style.accent]: accent,
                [style.row]: row,
                [style.fullWidth]: fullWidth,
                [style.noGap]: noGap,
                [style.center]: center,
            },
            className,
        )}
    >
        {children}
    </div>
)

Section.Title = Title
Section.Subtitle = Subtitle