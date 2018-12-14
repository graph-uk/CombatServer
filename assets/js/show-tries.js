const combat = window.combat = window.combat || {};
const NAVIGATION_LIST_CLASS = 'log-navigation';
const NAVIGATION_TITLE_CLASS = 'log-navigation__title';
const NAVIGATION_ITEM_CLASS = 'log-navigation__item';
const NAVIGATION_ITEM_ACTIVE_CLASS = 'log-navigation__item--active';
const LOG_CONTAINER_CLASS = 'log';

const renderTriesNavigation = (tries, lastSuccessfulRun, $tryPlaceholder, caseStatus) => {
	const {createTag, showTryDetails} = combat;
	const $buttons = tries.map((item, index) => {
		console.log(item);
		console.log(" is the item");
		const $el = createTag('div', {
			class: `${NAVIGATION_ITEM_CLASS}${index === 0 ? ' ' + NAVIGATION_ITEM_ACTIVE_CLASS : ''}`,
			children: 1 + index
		});

		$el.addEventListener('click', ({target}) => showTryDetails(item, $tryPlaceholder, target), false);

		return $el;
	});

    console.log($buttons);
    console.log(typeof $buttons);
    let lastSuccessfulRunElem =createTag('div', {
            class: `${NAVIGATION_ITEM_CLASS}`,
			children: "Last successful run"
        });
    if (lastSuccessfulRun.steps==null||lastSuccessfulRun.steps.length === 0||caseStatus === "success"){console.log("No successful runs"); currentSliderIndex = null;}else{
    	lastSuccessfulRunElem.addEventListener('click', ({target}) => showTryDetails(lastSuccessfulRun, $tryPlaceholder, target), false);
		$buttons.push(lastSuccessfulRunElem);
		console.log($buttons);

    }
	return {
		'$triesNavigation': createTag('div', {class: NAVIGATION_LIST_CLASS, children: [
			createTag('div', {class: NAVIGATION_TITLE_CLASS, children: 'Logs:'}),
			...$buttons
		]}),
		$buttons
	};
};

combat.showTries = ($el, tries, lastSuccessfulRun, caseStatus) => {
    console.log("Data in show tries");
    console.log(tries);
	const {createTag, showTryDetails} = combat;
	const $tryPlaceholder = createTag('div', {className: LOG_CONTAINER_CLASS});
	const {$triesNavigation, $buttons} = renderTriesNavigation(tries,lastSuccessfulRun, $tryPlaceholder, caseStatus);
	const stickyElem = createTag('div', {id: 'stickyElement'});
	stickyElem.append(
		$triesNavigation,
        $tryPlaceholder);
	$el.innerHTML = '';
	$el.append(
		stickyElem
	);

	showTryDetails(tries[0], $tryPlaceholder);
};
