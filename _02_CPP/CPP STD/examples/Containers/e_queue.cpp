#include <iostream>
#include <queue>
#include <string>
#include <deque>
#include <list>
#include <vector>
#include <functional>

module Containers;
using namespace std;

template<typename T, typename Container>
void print(const queue<T, Container>& q, const string& label) {
	cout << label << ": ";
	auto tmp = q;
	while (!tmp.empty()) {
		cout << tmp.front() << " ";
		tmp.pop();
	}
	cout << endl;
}
template<typename T, typename Container, typename Compare>
void print(const priority_queue<T, Container, Compare>& pq, const string& label) {
	cout << label << ": ";
	auto tmp = pq;
	while (!tmp.empty()) {
		cout << tmp.top() << " ";
		tmp.pop();
	}
	cout << endl;
}
static struct de_Person {
	string name;
	int age;
	de_Person(const string& n, int a) : name(n), age(a) {}
};
static ostream& operator<<(ostream& os, const de_Person& p) {
	return os << "{" << p.name << ", " << p.age << "}";
}

// std::queue
void example_queue() {
	cout << "\n=== std::queue ===\n";

	queue<int> q_int({ 1,2,3,4,5 });
	queue q_list{ list{ 10,20,30 } };
	queue q_deq{ deque{ 1,2,3,4,5,6,7,8,9 } };

	queue q(q_int);
	cout << "q size: " << q.size() << ", empty: " << boolalpha << q.empty() << endl;
	print(q, "q");

	q.front() = 999;
	q.back() = 666;
	print(q, "q (modified)");

	q.push_range(vector{ 6,7,8,9,10 });
	print(q, "q (push range)");

	// Queue with custom underlying container
	queue<string, list<string>> listQueue;
	listQueue.push("Custom");
	listQueue.push("Container");
	print(listQueue, "list-based queue");

	// Custom type in queue
	queue<de_Person> peopleQueue;
	peopleQueue.emplace("Alice", 30);
	peopleQueue.emplace("Bob", 25);
	peopleQueue.emplace("Charlie", 35);
	print(peopleQueue, "people queue");
}


// std::priority_queue
void example_priority_queue() {
	cout << "\n=== std::priority_queue ===\n";

	// less priority_queue
	priority_queue<int> d_pq;
	for (int i : {5, 2, 8, 1, 9}) d_pq.push(i);
	print(d_pq, "default Comp is less");

	auto vec = vector{ 90, 10, 80, 20, 70, 30, 60, 40, 50 };
	priority_queue pq_vec(greater<int>(), vec);
	priority_queue pq_deq(less<int>(), deque{ 100, 200, 300 });

	priority_queue pq(pq_vec);
	cout << "pq size: " << pq.size() << ", empty: " << boolalpha << pq.empty() << endl;
	print(pq, "pq");

	// push
	pq.push_range(list{ 1,30,500,7000,90000 });
	print(pq, "pq (push range)");

	// Remove top element
	pq.pop();
	print(pq, "pq (after pop)");


	// String priority queue
	priority_queue<string> stringPq;  // less is default Compare
	stringPq.push("Banana");
	stringPq.push("Apple");
	stringPq.push("Cherry");
	stringPq.push("Date");
	print(stringPq, "string priority queue");

	// Custom comparator
	struct PersonCompare {
		bool operator()(const de_Person& a, const de_Person& b) {
			return a.age > b.age; // ascending order by age
		}
	};
	priority_queue<de_Person, vector<de_Person>, PersonCompare> peopleQueue{ PersonCompare() };
	peopleQueue.emplace("Alice", 30);
	peopleQueue.emplace("Bob", 25);
	peopleQueue.emplace("Charlie", 35);
	print(peopleQueue, "people queue (sorted by age)");
}

void test_queue() {
	example_queue();
	example_priority_queue();
}

