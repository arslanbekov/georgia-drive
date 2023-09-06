new Vue({
    el: '#app',
    data: {
        loading: true,
        lastExecutionTime: '',
        activeRoute: '/api/get-auto',
        results: [],
        citiesOrder: [
            { name: "Kutaisi", order: 1 },
            { name: "Batumi", order: 2 },
            { name: "Telavi", order: 3 },
            { name: "Akhaltsikhe", order: 4 },
            { name: "Zugdidi", order: 5 },
            { name: "Gori", order: 6 },
            { name: "Poti", order: 7 },
            { name: "Ozurgeti", order: 8 },
            { name: "Sachkhere", order: 9 },
            { name: "Rustavi", order: 10 }
        ],
        currentRoute: '/api/get-auto',
        currentTitle: 'Drive city (auto)'
    },
    created() {
        this.fetchData();
        this.fetchData();
        this.fetchLastExecutionTime();
    },
    methods: {
        cityResults(cityName) {
            return this.results.filter(result => result.name === cityName);
        },
        isSaturday(dateString) {
            if (typeof dateString !== 'string') return false;
            const date = new Date(dateString);
            return date.getDay() === 6; // 6 - это суббота
        },
        setCurrentRoute(route) {
            this.loading = true;
            this.activeRoute = route;
            this.currentRoute = route;

            if (route === '/api/get-auto') {
                this.currentTitle = 'Drive city (auto)';
            } else if (route === '/api/get-manual') {
                this.currentTitle = 'Drive city (manual)';
            } else if (route === '/api/get-theory') {
                this.currentTitle = 'Theory';
            }

            this.fetchData();
        },
        fetchLastExecutionTime() {
            axios.get("/api/last-exec-time")
                .then(response => {
                    const timestamp = response.data[0].timestamp;
                    const dateObj = new Date(timestamp);
                    this.lastExecutionTime = this.formatDate(dateObj);
                })
                .catch(error => {
                    console.error("Error fetching last execution time:", error);
                });
        },
        formatDate(dateObj) {
            // Форматируем дату в человекочитаемый формат
            return `${dateObj.toLocaleDateString()} ${dateObj.toLocaleTimeString()}`;
        },
        fetchData() {
            fetch(this.currentRoute)
                .then(response => response.json())
                .then(data => {
                    this.results = data;
                    this.loading = false;
                })
                .catch(error => {
                    console.error("There was an error fetching the data:", error);
                    this.loading = false;
                });
        }
    }
});
