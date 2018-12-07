const combat = window.combat = window.combat || {};
const NAVIGATION_ITEM_ACTIVE_CLASS = 'log-navigation__item--active';
const LOG_TABS_CLASS = 'log__tabs';
const LOG_NAV_ITEM_CLASS = 'log__tabs-nav';
const LOG_NAV_ACTIVE_CLASS = 'log__tabs-nav--active';
const LOG_SLIDER_CLASS = 'log__details';
const LOG_TERMINAL_CLASS = 'log__terminal';

combat.showTryDetails = (tryData, $placeholder, $target) => {
	if (!$target) {
		return renderTry($placeholder, tryData);
	}

	const $prevActive = $target.parentElement.querySelector(`.${NAVIGATION_ITEM_ACTIVE_CLASS}`);

	if ($target !== $prevActive) {
		if ($prevActive) {
			$prevActive.className = $prevActive.className.replace(` ${NAVIGATION_ITEM_ACTIVE_CLASS}`, '');
		}

		$target.className += ` ${NAVIGATION_ITEM_ACTIVE_CLASS}`;

		renderTry($placeholder, tryData);
	}
}

const renderTry = ($placeholder, data) => {
	const {createTag} = combat;
	const $detailsPlaceholder = combat.createTag('div', {class: LOG_SLIDER_CLASS});
	const $buttons = renderTabs(data, $detailsPlaceholder, key => onTabClick($buttons, key));

	$placeholder.innerHTML = '';
	$placeholder.append(
		$buttons,
		$detailsPlaceholder
	);

	const firstKey = Object.keys(data)[0];

	renderDetails(data[firstKey], $detailsPlaceholder);
	onTabClick($buttons, firstKey);
}

const onTabClick = ($buttons, dataKey) => {
	const $lastActive = $buttons.querySelector(`.${LOG_NAV_ACTIVE_CLASS}`);
	const classString = ` ${LOG_NAV_ACTIVE_CLASS}`;

	if ($lastActive) {
		$lastActive.className = $lastActive.className.replace(classString, '');
	}
	$buttons.querySelector(`[data-key=${dataKey}]`).className += classString;
};

const renderTabs = (data, $detailsPlaceholder, onTabClick) => {
	const {createTag} = combat;
	const $buttons = Object.keys(data).map(key => {
		const $btn = createTag('div', {
			class: `col ${LOG_NAV_ITEM_CLASS}`,
			'data-key': key,
			children: key
		});

		$btn.addEventListener('click', () => {
			renderDetails(data[key], $detailsPlaceholder);
			onTabClick(key);
		}, false);

		return $btn;
	});

	return createTag('div', {class: LOG_TABS_CLASS, children:
		createTag('div', {class: 'row no-gutters', children:
			$buttons
		})
	});
}

const renderDetails = (data, $placeholder) => {
	const {createTag, renderSlider} = combat;

	$placeholder.innerHTML = '';

	if (typeof data === 'string') {
		return $placeholder.append(
			createTag('div', {class: LOG_TERMINAL_CLASS, children: data})
		);
	}

	if (data && data.constructor === Array) {
		return renderSlider(data, $placeholder);
	}
}
