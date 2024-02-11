import clsx from 'clsx';
import Heading from '@theme/Heading';
import styles from './styles.module.css';

const FeatureList = [
  {
    title: 'Photobooth',
    Svg: require('@site/static/img/illu_camera.svg').default,
    description: (
      <>
        Enhance your parties and keep memories by bringing your home-made photobooth.
      </>
    ),
  },
  {
    title: 'Karaoke',
    Svg: require('@site/static/img/illu_karaoke.svg').default,
    description: (
      <>
        Sing along to your favorite song, supporting MP3+CDG and videos.
      </>
    ),
  },
  {
    title: 'Home-made',
    Svg: require('@site/static/img/illu_diy.svg').default,
    description: (
      <>
        Build your own appliance and bring the hardware that matches your budget. Works from Raspberry Pis to full desktop computers.
      </>
    ),
  },
];

function Feature({Svg, title, description}) {
  return (
    <div className={clsx('col col--4')}>
      <div className="text--center">
        <Svg className={styles.featureSvg} role="img" />
      </div>
      <div className="text--center padding-horiz--md">
        <Heading as="h3">{title}</Heading>
        <p>{description}</p>
      </div>
    </div>
  );
}

export default function HomepageFeatures() {
  return (
    <section className={styles.features}>
      <div className="container">
        <div className="row">
          {FeatureList.map((props, idx) => (
            <Feature key={idx} {...props} />
          ))}
        </div>
      </div>
    </section>
  );
}
