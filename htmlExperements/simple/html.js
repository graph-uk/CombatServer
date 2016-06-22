function Spoil(tryID) {
	//alert("sdf");
	if (document.getElementById(tryID).style.display != '') {
	    document.getElementById(tryID).style.display = '';
	    this.innerText = '';
	    this.value = 'Minimize';
	} else {
	    document.getElementById(tryID).style.display = 'none';
	    this.innerText = '';
	    this.value = 'Try1';
	}
}