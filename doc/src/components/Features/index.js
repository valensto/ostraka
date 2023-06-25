import React from 'react';
import clsx from 'clsx';
import styles from './styles.module.css';
import {Section} from "../Section";
import Link from "@docusaurus/Link";

const FeatureList = [
    {
        title: 'Simplified Configuration',
        points: [
            "Intuitive YAML setup",
            "Easy field mapping",
            "Customizable middlewares",
        ],
        ctaLink: "",
        ctaTitle: "Lean Configuration"
    },
    {
        title: 'Flexible Event Routing',
        points: [
            "Versatile integration options",
            "Diverse destination choices",
            "Configurable routing conditions",
        ],
        ctaLink: "",
        ctaTitle: "Lean How to Route"
    },
    {
        title: 'Event Processing Efficiency',
        points: [
            "High-performance event handling",
            "Parallel processing capabilities",
            "Efficient resource utilization",
        ],
        ctaLink: "",
        ctaTitle: "See Performance"
    },
];

function Feature({title, points, ctaLink, ctaTitle}) {
    return (
        <div className='col col--4'>
            <div className={clsx("padding-horiz--md", styles.feature)}>
                <h3 className={"text--center"}>{title}</h3>
                <ul className={styles.list}>
                    {points.map((point, idx) => (
                        <li className={styles.listItem} key={idx}>{point}</li>
                    ))}
                </ul>

                <div className={styles.cta}>
                    <div>
                        <Link
                            className={clsx('button button--lg', styles.btnGradient)}
                            to="/docs/intro">
                            {ctaTitle}
                        </Link>
                    </div>
                </div>
            </div>
        </div>
    );
}

export default function HomepageFeatures() {
    return (
        <Section fullWidth center>
            <div className="container">
                <div className="row">
                    {FeatureList.map((props, idx) => (
                        <Feature key={idx} {...props} />
                    ))}
                </div>
            </div>
        </Section>
    );
}
