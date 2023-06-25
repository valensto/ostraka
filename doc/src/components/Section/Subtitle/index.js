import React from "react"
import style from "./styles.module.css"
import clsx from "clsx"

export const Subtitle = ({
                             children,
                             center,
                             size = "medium",
                             className = "",
                         }) => (
    <p
        className={clsx(
            style.root,
            { [style.center]: center },
            style[`size-${size}`],
            className,
        )}
    >
        {children}
    </p>
)