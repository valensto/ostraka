import React from 'react';
import clsx from 'clsx';
import Link from '@docusaurus/Link';
import useDocusaurusContext from '@docusaurus/useDocusaurusContext';
import Layout from '@theme/Layout';
import Features from '@site/src/components/Features';
import {Builder} from '@site/src/components/Builder';

import styles from './index.module.css';
import {WhyCards} from "../components/Why";

function HomepageHeader() {
    const {siteConfig} = useDocusaurusContext();
    return (
        <header className={clsx('hero', styles.heroBanner)}>
            <div className="container">
                <h1 className="hero__title">{siteConfig.tagline}</h1>
                <p className={clsx('hero__subtitle', styles.description)}>
                    Route events from diverse sources to multiple destinations with Ostraka, an open-source event
                    dispatching tool.
                </p>
                <div className={styles.buttons}>
                    <Link
                        className={clsx('button button--lg', styles.btnGradient)}
                        to="/docs/getting-started">
                        Dispatch your events - 5min ⏱️
                    </Link>
                </div>
            </div>
        </header>
    );
}

export default function Home() {
    return (
        <Layout
            title={`Effortless event dispatcher`}
            description="Route events from diverse sources to multiple destinations with Ostraka, an open-source event dispatching tool.">
            <HomepageHeader/>
            <main>
                <Features/>
                <Builder/>
                <WhyCards/>
            </main>
        </Layout>
    );
}
