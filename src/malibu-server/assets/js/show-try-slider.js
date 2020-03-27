const combat = window.combat = window.combat || {};
const LOG_SLIDE_ID = 'log_slider';
const LOG_SLIDE_CLASS = 'log__slide';
const LOG_SLIDE_URL_CLASS = 'log__slide-url';
const LOG_SLIDE_SOURCE_CLASS = 'log__slide-source';
const LOG_SLIDE_IMG_CLASS = 'log__slide-img';
const LOG_SLIDE_IMG_LINK_CLASS = 'log__slide-img-link';

combat.renderSlider = (data, $target) => {
	const $slider = createSliderMarkUp(data);

	$target.append($slider);

	if (combat._sliderInstace) {
		combat._sliderInstace.destroy();
		combat._sliderInstace = undefined;
	}

	if (combat._sliderInterval) {
		clearInterval(combat._sliderInterval);
		combat._sliderInterval = undefined;
	}

	// if (combat.slides.all > 0) {
	// 	const check = () => {
	// 		if (combat.slides.all === combat.slides.loaded) {
	// 			clearInterval(combat._sliderInterval);
	// 			combat._sliderInterval = undefined;
	// 			combat._sliderInstace = new Glide($slider).mount();
	// 		}
	// 	}
	//
	// 	combat._sliderInterval = setInterval(check, 10);
	// 	check();
	// } else {
		combat._sliderInstace = new Glide($slider, { startAt: currentSliderIndex==null? combat.slides.all-1:currentSliderIndex<=combat.slides.all-1?currentSliderIndex:combat.slides.all-1 }).mount();
	// }
};

const createSliderMarkUp = data => {
	const {createTag} = combat;
	const $slides = createTag('ul', {class: 'glide__slides'});
	const $bullets = createTag('div', {class: 'slider__bullets glide__bullets', 'data-glide-el': 'controls[nav]'});

	combat.slides = {all: 0, loaded: 0};

	data.reverse();
	data.forEach(({image, source, url}, index) => {
        index = Math.abs(data.length-index);
		const $slide = createTag('li', {class: `glide__slide ${LOG_SLIDE_CLASS}`});

		if (url) {
			$slide.append(createTag('a', {
				class: LOG_SLIDE_URL_CLASS,
				href: url,
				target: '_blank',
				children: '> log url'
			}));
		}

		if (source) {
			$slide.append(createTag('a', {
				class: LOG_SLIDE_SOURCE_CLASS,
				href: source,
				target: '_blank',
				children: '> log source'
			}));
		}

		if (image) {
			combat.slides.all += 1;

			const $img = createTag('img', {
				class: LOG_SLIDE_IMG_CLASS
			});

			$img.addEventListener('load', () => combat.slides.loaded += 1, false);
			// $img.addEventListener('load', () => console.log(index+ " out of " + combat.slides.all))
			$img.setAttribute('src', image);

			$slide.append(createTag('a', {
				class: LOG_SLIDE_IMG_LINK_CLASS,
				href: image,
				target: '_blank',
				children: $img,
				index: (index-1)
			}));
		}

		$slides.prepend($slide);
		$bullets.prepend(createTag('button', {'class': 'slider__bullet glide__bullet', 'data-glide-dir': "="+(index-1)}));
	});
	data.reverse();
	if(currentSliderIndex>combat.slides.all-1||currentSliderIndex===null) {
        currentSliderIndex = combat.slides.all-1;
    }
    console.log("Current slider index is " + currentSliderIndex);
	return createTag('div', {class: 'glide', id: LOG_SLIDE_ID, children: [
		createTag('div', {class: 'glide__track', 'data-glide-el': 'track', children:
			$slides
		}),
		createTag('div', {class: 'glide__arrows', 'data-glide-el': 'controls', children: [
			createTag('button', {class: 'glide__arrow glide__arrow--left',
				'data-glide-dir': '<', children: '<'
			}),
			createTag('button', { class: 'glide__arrow glide__arrow--right',
				'data-glide-dir': '>', children: '>'
			})
		]})
		,$bullets,
            createTag('span', { id: 'test_number_counter', children : currentSliderIndex+1 + ' out of ' + combat.slides.all
            })
	]});
};
