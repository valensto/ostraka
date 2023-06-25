import { cloneElement } from "react"

const SvgImage = ({ image, title = "" }) =>
    cloneElement(image, {
        ...image.props,
        title,
    })

export default SvgImage
